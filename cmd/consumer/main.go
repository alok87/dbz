package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"k8s.io/klog/v2"

	"github.com/Shopify/sarama"
	"github.com/spf13/pflag"
)

func init() {
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))
}

func main() {
	klog.Infof("starting consumer...")
	time.Sleep(time.Second * 5)

	klog.Infof("started consumer")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up configuration
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Connect to Kafka broker
	brokers := []string{"kafka:9092"}
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		klog.Fatalf("Failed to start consumer:", err)
	}
	defer consumer.Close()

	// Define the topic you want to consume from
	topic := "dbserver1.inventory.customers"
	partitionConsumer, err := consumer.ConsumePartition(
		topic, 0, sarama.OffsetNewest)
	if err != nil {
		klog.Fatalf("Failed to start partition consumer:", err)
	}
	defer partitionConsumer.Close()

	wg := &sync.WaitGroup{}

	klog.Infof("run consumer")

	// Start a goroutine to consume messages
	wg.Add(1)
	go func() {
		defer wg.Done()
		klog.Infof("consumer started, waiting for msgs")
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				klog.Infof("received message")
				_ = msg
			case err := <-partitionConsumer.Errors():
				klog.Infof("error consuming messages:", err)
			case <-ctx.Done():
				klog.Infof("exiting consumer, ctx done")
				return
			}
		}
	}()

	// graceful shutdown, no libs required, understand just below
	wg.Add(1)
	go func() {
		defer wg.Done()

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigterm:
			klog.Infof("sigterm received")
		case <-ctx.Done():
			klog.Infof("context done, bye")
			return
		}

		cancel()
	}()

	wg.Wait()
	klog.Infof("shutdown")
}
