package main

import (
	"fmt"
	"github.com/hi-devops-io/hi-devops/pkg/config"
)

func main()  {
	fmt.Println("hi devops")
	javaConfig := config.BuildConfig("dev", "demo", "my-app", "1.0.0","java")
	fmt.Println(javaConfig)

	javaWarConfig := config.BuildConfig("dev", "demo", "my-app","1.0.0","java-war")
	fmt.Println(javaWarConfig)
}
