package domain

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig *string
var ron = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(i int) int {
	return ron.Intn(i)
}

type PodManage struct {
	Ctx     context.Context
	Cancel  context.CancelFunc
	Started bool
}

func Connect() (*kubernetes.Clientset, error) {
	var config, err = rest.InClusterConfig()
	if err == nil {
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return clientset, nil
	} else {
		if kubeconfig == nil {
			if home := homeDir(); home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()
		}

		log.Println("Running out of Kubernetes cluster")
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
			return nil, err
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
			return nil, err
		}
		return clientset, nil
	}
}

func (m *PodManage) Start(c *Config) error {
	log.Println("Cookie Time!!! Random feast starting")

	err := m.MainLoop(c)
	if err != nil {
		return err
	}

	timeout := time.After(time.Duration(c.Duration) * time.Second)
	go func() {
		tick := time.Tick(time.Duration(c.Interval) * time.Second)
		for {
			select {
			case <-timeout:
				log.Println("Cookie Monster Duration Timeout !!!")
				m.Stop(c)
			case <-m.Ctx.Done():
				return
			case <-tick:
				m.MainLoop(c)
			}
		}
	}()

	return nil
}

func (m *PodManage) MainLoop(c *Config) error {
	for _, ns := range c.Namespace {
		for _, res := range ns.Resource {
			for nu := 0; nu < int(res.Target); nu++ {
				pod, startKill, err := m.SelectVictimPod(c, ns.Name, res.Kind, res.Name)
				if err != nil {
					log.Println(err)
					return err
				} else if startKill {
					go killPod(pod, ns.Name)
				}
			}
		}
	}
	return nil
}

func (m *PodManage) Stop(c *Config) {
	log.Println("Stop snacking.\n")

	defer m.Cancel()
	m.Started = false
}

func (m *PodManage) SelectVictimPod(c *Config, ns string, kind string, name string) (*v1.Pod, bool, error) {
	con, err := Connect()
	if err != nil {
		log.Println(err)
		return nil, false, err
	}

	var matchLabels map[string]string
	var lo metav1.ListOptions
	if name == "" {
		lo = metav1.ListOptions{}
	} else {
		lo = metav1.ListOptions{
			FieldSelector: "metadata.name=" + name,
		}
	}

	switch kind {
	case "deployment":
		deploymentsClient := con.AppsV1().Deployments(ns)

		deps, err := deploymentsClient.List(lo)
		if err != nil {
			fmt.Println(err)
			return nil, false, err
		} else if len(deps.Items) < 1 {
			log.Printf("Can not find %s %s in namespace %s, doing nothing", kind, name, ns)
			return nil, false, err
		}

		x := RandomInt(len(deps.Items))
		dep := deps.Items[x]
		fmt.Printf(" * %s (%d replicas)\n", dep.Name, *dep.Spec.Replicas)
		log.Printf("%s %s in namespace %s has %d pods defined, %d available and %d unavailable ",
			kind, dep.ObjectMeta.Name, ns, *dep.Spec.Replicas, dep.Status.AvailableReplicas, dep.Status.UnavailableReplicas)

		if dep.ObjectMeta.Name == "rabbitmq" {
			if (*dep.Spec.Replicas/2 + 1) >= dep.Status.AvailableReplicas {
				log.Printf("available pods less than unavailable pods, doing nothing")
				return nil, false, nil
			}
		} else {
			if dep.Status.AvailableReplicas < 2 {
				log.Printf("Only one pod available, doing nothing")
				return nil, false, nil
			}
		}

		matchLabels = dep.Spec.Selector.MatchLabels
	case "statefulset":
		statefulClient := con.AppsV1().StatefulSets(ns)
		sss, err := statefulClient.List(lo)
		if err != nil {
			fmt.Println(err)
			return nil, false, err
		} else if len(sss.Items) < 1 {
			log.Printf("Can not find %s %s in namespace %s, doing nothing", kind, name, ns)
			return nil, false, err
		}

		x := RandomInt(len(sss.Items))
		ss := sss.Items[x]
		fmt.Printf(" * %s (%d replicas)\n", ss.Name, *ss.Spec.Replicas)
		log.Printf("%s %s in namespace %s has %d pods defined, %d available ",
			kind, ss.ObjectMeta.Name, ns, *ss.Spec.Replicas, ss.Status.ReadyReplicas)

		if ss.ObjectMeta.Name == "mariadb" {
			if (*ss.Spec.Replicas/2 + 1) >= ss.Status.ReadyReplicas {
				log.Printf("available pods less than unavailable pods, doing nothing")
				return nil, false, nil
			}
		} else {
			if ss.Status.ReadyReplicas < 2 {
				log.Printf("Only one pod available, doing nothing")
				return nil, false, nil
			}
		}

		matchLabels = ss.Spec.Selector.MatchLabels
	case "daemonset":
		daemonsetClient := con.AppsV1().Deployments(ns)

		dss, err := daemonsetClient.List(lo)
		if err != nil {
			fmt.Println(err)
			return nil, false, err
		} else if len(dss.Items) < 1 {
			log.Printf("Can not find %s %s in namespace %s, doing nothing", kind, name, ns)
			return nil, false, err
		}

		x := RandomInt(len(dss.Items))
		ds := dss.Items[x]
		fmt.Printf(" * %s (%d replicas)\n", ds.Name, *ds.Spec.Replicas)
		log.Printf("%s %s in namespace %s has %d pods defined, %d available and %d unavailable ",
			kind, ds.ObjectMeta.Name, ns, *ds.Spec.Replicas, ds.Status.AvailableReplicas, ds.Status.UnavailableReplicas)

		if ds.Status.AvailableReplicas < 2 {
			log.Printf("Only one pod available, doing nothing")
			return nil, false, nil
		}

		matchLabels = ds.Spec.Selector.MatchLabels
	}

	s := ""
	count := 0
	for k, v := range matchLabels {
		count += 1
		s = s + k + "=" + v
		if count < len(matchLabels) {
			s = s + ","
		}
	}
	lo = metav1.ListOptions{
		FieldSelector: "status.phase=Running",
		LabelSelector: s,
	}

	pods, err := con.CoreV1().Pods(ns).List(lo)
	if err != nil {
		fmt.Println(err)
		return nil, false, err
	}
	x := RandomInt(len(pods.Items))
	pod := pods.Items[x]

	return &pod, true, nil
}

func killPod(pod *v1.Pod, ns string) error {
	con, err := Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	err = con.CoreV1().Pods(ns).Delete(pod.Name, &metav1.DeleteOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Eating pod %s NOM NOM NOM!!!!", pod.Name)

	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
