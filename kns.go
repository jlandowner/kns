package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Kubeconfig is a kubeconfig yaml struct
type Kubeconfig struct {
	APIVersion     string `yaml:"apiVersion"`
	Kind           string
	Clusters       []interface{}
	Contexts       []Context
	CurrentContext string `yaml:"current-context"`
	Preferences    interface{}
	Users          interface{}
}

// Context is a kubeconfig context struct
type Context struct {
	Context struct {
		Cluster   string
		Namespace string
		User      string
	}
	Name string
}

func main() {
	// flag check
	var kubeconfigPath *string
	kubeconfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	if *kubeconfigPath == "" {
		*kubeconfigPath = filepath.Join(homeDir(), ".kube", "config")
	}

	// Param check
	namespace := ""
	if len(flag.Args()) == 1 {
		switch flag.Args()[0] {
		case "kube-system", "kube", "sys", "system":
			namespace = "kube-system"
		case "default", "reset":
			namespace = "default"
		case "help":
			usage()
			return
		case "version":
			version()
			return
		default:
			namespace = flag.Args()[0]
		}
	} else if len(flag.Args()) > 1 {
		usage()
		return
	}
	// fmt.Println(namespace)

	// if namespace is default then skip cluster access process
	// it is based on the idea "default namespace exists in almost all clusters"
	if namespace != "default" {
		// Get namespaces from cluster of current contexts
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		// get Namespaces
		ns, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		// param as Namespace number
		if i, err := strconv.Atoi(namespace); err == nil {
			if i < len(ns.Items) {
				namespace = ns.Items[i].ObjectMeta.Name
			} else {
				namespace = ""
			}
		}

		// if no args then show namespaces list
		// if arg exists then validate namespaces
		if namespace == "" {
			fmt.Println("** List of Namespaces in the Current-context Cluster.")
			for i, v := range ns.Items {
				fmt.Println(i, ": ", v.ObjectMeta.Name)
			}
			if i, err := askNamespaceNum(len(ns.Items)); err != nil {
				fmt.Println("Exit.")
				return
			} else {
				namespace = ns.Items[i].ObjectMeta.Name
			}
		} else {
			isNsInCluster := false
			for _, v := range ns.Items {
				if namespace == v.ObjectMeta.Name {
					isNsInCluster = true
				}
			}
			if !isNsInCluster {
				fmt.Printf("Namespace %s does NOT Exist in the Cluster.\n", namespace)
				return
			}
		}
	}

	// Open kubeconfig File
	kConfig := readKubeconfig(*kubeconfigPath)

	// Matching Contexts
	currentContext := kConfig.CurrentContext
	numContext := 0
	isMatched := false
	for i, v := range kConfig.Contexts {
		if currentContext == v.Name {
			numContext = i
			isMatched = true
			break
		}
	}
	// Set namespace to current-context
	if isMatched {
		kConfig.Contexts[numContext].Context.Namespace = namespace
	} else {
		fmt.Println("Not Matched Context")
		return
	}

	// Write kubeconfig
	if err := writeKubeconfig(*kubeconfigPath, kConfig); err != nil {
		panic(err.Error())
	}

	fmt.Println("** Completed: Switch namespace ", namespace)
}

// get homeDir from env
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// read kubeconfig file
func readKubeconfig(kubeconfigPath string) (kConfig Kubeconfig) {
	buf, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		panic(err)
	}

	// yaml to struct
	err = yaml.Unmarshal(buf, &kConfig)
	if err != nil {
		panic(err)
	}
	return kConfig
}

// write kubeconfig to file
func writeKubeconfig(kubeconfigPath string, kConfig Kubeconfig) (err error) {
	err = nil

	// struct to yaml
	out, err := yaml.Marshal(kConfig)
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
	// Write kubeconfig file
	err = ioutil.WriteFile(kubeconfigPath, out, 0600)
	if err != nil {
		log.Fatalf("ERR: %v", err)
	}
	return err
}

// show dialog to ask
func askNamespaceNum(max int) (i int, qerr error) {
	fmt.Print("** Which namespace do you want to switch? (exit: q)\n")
	fmt.Print("Select[n] => ")
	i = 99
	qerr = nil
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		if atoi, err := strconv.Atoi(input); err == nil {
			i = atoi
			if 0 <= i && i < max {
				break
			}
			fmt.Print("Select[n] => ")
		} else if input == "q" {
			qerr = errors.New("quit")
			break
		} else {
			fmt.Print("Select[n] => ")
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return i, qerr
}

// how to use
func usage() {
	fmt.Fprintf(os.Stderr, "usage: kns [--kubeconfig path] [namespace]\n")
	flag.PrintDefaults()
}

func version() {
	version := "v0.1.0"
	auther := "jlandowner"
	repo := "https://github.com/jlandowner/kns"
	date := "2019/10/01 JST"
	fmt.Printf("version: %s, Auther: %s, Github: %s, Date: %s\n", version, auther, repo, date)
}
