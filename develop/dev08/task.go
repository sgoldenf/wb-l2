package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	// Использовала сторонний пакет, т.к. у меня не linux устройство (macOS, нет /proc)
	gops "github.com/mitchellh/go-ps"
)

/*
Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/

func main() {
	s := bufio.NewScanner(os.Stdin)
	for {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			continue
		}
		base := filepath.Base(path)
		fmt.Printf("%s$ ", base)
		s.Scan()
		osArgs := strings.Fields(s.Text())
		if len(osArgs) == 0 {
			continue
		}
		if f, ok := funcMap[osArgs[0]]; !ok {
			forkExec(osArgs)
		} else {
			f(osArgs)
		}
	}
}

var funcMap = map[string]func(args []string){
	"cd":   cd,
	"pwd":  pwd,
	"echo": echo,
	"kill": kill,
	"ps":   ps,
	"exit": exit,
}

func cd(args []string) {
	switch len(args) {
	case 1:
		if home, err := os.UserHomeDir(); err != nil {
			fmt.Println(err)
		} else {
			if err := os.Chdir(home); err != nil {
				fmt.Println(err)
			}
		}
	case 2:
		if err := os.Chdir(args[1]); err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("Too many args for cd command")
	}
}

func pwd(args []string) {
	if len(args) > 1 {
		fmt.Println("pwd: expected 0 arguments; got", len(args)-1)
	} else {
		if path, err := os.Getwd(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(path)
		}
	}
}

func echo(args []string) {
	fmt.Println(strings.Join(args[1:], " "))
}

func kill(args []string) {
	for _, pidSTR := range args[1:] {
		pid, err := strconv.Atoi(pidSTR)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func ps(args []string) {
	if len(args) > 1 {
		fmt.Println("Too many arguments")
		return
	}
	processes, err := gops.Processes()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%6s\t%s\n", "PID", "CMD")
	for _, process := range processes {
		fmt.Printf("%6d\t%s\n", process.Pid(), process.Executable())
	}
}

func forkExec(args []string) {
	path, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	args[0] = path
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	p, err := os.StartProcess(args[0], args, &procAttr)
	if err != nil {
		fmt.Println(err)
		return
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigCh:
		return
	default:
		state, err := p.Wait()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("exit(%d)\n", state.ExitCode())
		}
	}
}

func exit(args []string) {
	fmt.Println("exit")
	os.Exit(0)
}
