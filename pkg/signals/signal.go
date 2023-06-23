/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package signals

import (
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		// 在接收到第一个信号后，关闭 stop 通道，表示程序应该开始进行优雅地关闭（例如，等待正在进行的工作完成，清理资源等）
		close(stop)
		// 在接收到第二个信号后，调用 os.Exit(1) 直接终止程序，而不进行任何优雅地关闭。
		//这种情况通常发生在用户希望立即停止程序时（例如，连续发送两次信号）
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
