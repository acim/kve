package main

import (
	"encoding/json"
	"log"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type client struct {
	kubernetes.Interface
}

func newClient(clientset kubernetes.Interface) *client {
	return &client{Interface: clientset}
}

func (c *client) classes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	scs, err := c.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		httpError(w, err, "failed getting storage classes")

		return
	}

	for i := range scs.Items {
		scs.Items[i].ManagedFields = nil
	}

	resp := resClasses{
		StorageClasses: scs.Items,
	}

	res, err := json.Marshal(resp)
	if err != nil {
		httpError(w, err, "failed encoding to json")

		return
	}

	w.Write(res) //nolint:errcheck
}

func (c *client) volumes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	pvs, err := c.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		httpError(w, err, "failed getting persistent volumes")

		return
	}

	pvcs, err := c.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil {
		httpError(w, err, "failed getting persistent volume claims")

		return
	}

	pods, err := c.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		httpError(w, err, "failed getting pods")

		return
	}

	// FieldSelector: fields.Set{"spec.volumes[].persistentVolumeClaim.claimName": "ghost-acim"}.AsSelector().String(),
	// _ = corev1.ReadWriteMany

	resp := resVolumes{
		Volumes: getVolumes(pvs.Items, pvcs.Items, pods.Items),
	}

	res, err := json.Marshal(resp)
	if err != nil {
		httpError(w, err, "failed encoding to json")

		return
	}

	w.Write(res) //nolint:errcheck
}

func httpError(w http.ResponseWriter, err error, text string) {
	log.Printf("%s: %v\n", text, err)

	r := resError{
		Error: &text,
	}

	res, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(res) //nolint:errcheck
}

func getVolumes(pvs []corev1.PersistentVolume, pvcs []corev1.PersistentVolumeClaim, pods []corev1.Pod) []volume {
	volumes := make([]volume, len(pvs))

	for i, pv := range pvs {
		pv.ManagedFields = nil
		volumes[i].PersistentVolume = pv

		for _, pvc := range pvcs {
			if pv.Name == pvc.Spec.VolumeName {
				pvc.ManagedFields = nil
				volumes[i].PersistentVolumeClaim = pvc

				for _, pod := range pods {
					for _, v := range pod.Spec.Volumes {
						if v.PersistentVolumeClaim != nil && pvc.Name == v.PersistentVolumeClaim.ClaimName {
							pod.ManagedFields = nil
							volumes[i].Pods = append(volumes[i].Pods, pod)

							break
						}
					}
				}

				break
			}
		}
	}

	return volumes
}

type resClasses struct {
	StorageClasses []storagev1.StorageClass `json:"classes"`
}

type resVolumes struct {
	Volumes []volume `json:"volumes"`
}

type resError struct {
	Error *string `json:"error,omitempty"`
}

type volume struct {
	PersistentVolume      corev1.PersistentVolume      `json:"volume"`
	PersistentVolumeClaim corev1.PersistentVolumeClaim `json:"claim"`
	Pods                  []corev1.Pod                 `json:"pods"`
}
