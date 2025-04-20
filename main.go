package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/rumpl/devoxx-docker/remote"
)

func main() {
	//log.Println("Args:", os.Args)
	if len(os.Args) < 2 {
		if err := run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
	switch os.Args[1] {
	case "child":
		if len(os.Args) < 3 {
			log.Fatal("Missing image name")
		}
		if err := child(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	case "pull":
		if len(os.Args) < 3 {
			log.Fatal("Missing image name")
		}
		if err := pull(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	case "run":
		if len(os.Args) < 4 {
			log.Fatal("Missing image name or command to run")
		}
		if err := run(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unknown command", os.Args[1])
	}
}
func pull(image string) error {
	fmt.Printf("Pulling %s\n", image)
	puller := remote.NewImagePuller(image)
	if err := puller.Pull(); err != nil {
		return fmt.Errorf("pull failed: %w", err)
	}
	fmt.Println("Pulling done")
	return nil
}

func child(image string) error {
	fmt.Printf("CHILD PID: %d\n", os.Getpid())

	if err := syscall.Sethostname([]byte("container")); err != nil {
		return fmt.Errorf("sethostname failed: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	fmt.Printf("CHILD Hostname: %s\n", hostname)

	// Change root directory
	if err := syscall.Chroot(fmt.Sprintf("/fs/%s", image)); err != nil {
		return fmt.Errorf("chroot failed: %w", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir failed: %w", err)
	}

	// Mount proc filesystem
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("mount proc failed: %w", err)
	}

	// Execute the command
	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func run() error {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Add mount namespace along with existing namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start failed: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("wait failed: %w", err)
	}

	fmt.Printf("Container exited with code %d\n", cmd.ProcessState.ExitCode())
	return nil
}
