package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cnrancher/autok3s/pkg/cluster"
	"github.com/cnrancher/autok3s/pkg/common"
	"github.com/cnrancher/autok3s/pkg/providers"
	"github.com/cnrancher/autok3s/pkg/types"
	autok3stypes "github.com/cnrancher/autok3s/pkg/types/apis"
	"github.com/cnrancher/autok3s/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//call comamnd kubectl drain nodename to drain all pods of nodename firstly then call kubectl delete node nodename to delete the node from cluster
//added by edgego
func DeleteCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster"]
	provider := vars["provider"]

	if clusterID == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("clusterID cannot be empty"))
		return
	}

	if provider == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("provider cannot be empty"))
		return
	}

	state, err := common.DefaultDB.GetClusterByID(clusterID)
	if err != nil || state == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("cluster %s is not found", clusterID)))
		return
	}

	p, err := providers.GetProvider(state.Provider)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("provider %s is not found", state.Provider)))
		return
	}

	opt, err := p.GetProviderOptions(state.Options)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("provider option%s is not found", state.Provider)))
		return
	}

	cluster := &autok3stypes.Cluster{
		Metadata: state.Metadata,
		Options:  opt,
	}

	b, err := json.Marshal(cluster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Marshal cluster: %#v failed ,error: %s", cluster, err.Error())))
		return
	}

	err = p.SetConfig(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to call SetConfig,error: %s", err.Error())))
		return
	}

	err = p.MergeClusterOptions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to call MergeClusterOptions,error: %s", err.Error())))
		return
	}

	p.GenerateClusterName()

	err = p.DeleteK3sCluster(true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to call DeleteK3sCluster,error: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster"]
	nodeName := vars["node"]
	instanceId := vars["instance"]

	if clusterID == "" {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("DeleteK3sNode: argument cluster : [%s] is empty", clusterID)))
		return
	}

	if nodeName == "" {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("DeleteK3sNode: argument node : [%s] is empty", nodeName)))
		return
	}

	if instanceId == "" {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("DeleteK3sNode: argument instanceId : [%s] is empty", instanceId)))
		return
	}

	state, err := common.DefaultDB.GetClusterByID(clusterID)
	if err != nil || state == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("cluster %s is not found", clusterID)))
		return
	}

	nodeIp := strings.Replace(instanceId, "-", ".", -1)
	if nodeIp == state.IP {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("DeleteK3sNode: can not delete first node, ip : [%s],host name :[%s]", nodeIp, nodeName)))
		return
	}

	p, err := providers.GetProvider(state.Provider)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("provider %s is not found", state.Provider)))
		return
	}

	opt, err := p.GetProviderOptions(state.Options)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("provider option%s is not found", state.Provider)))
		return
	}

	cluster := &autok3stypes.Cluster{
		Metadata: state.Metadata,
		Options:  opt,
	}

	b, err := json.Marshal(cluster)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Marshal cluster: %#v failed ,error: %s", cluster, err.Error())))
		return
	}

	err = p.SetConfig(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to call SetConfig,error: %s", err.Error())))
		return
	}

	err = p.MergeClusterOptions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to call MergeClusterOptions,error: %s", err.Error())))
		return
	}

	id := p.GenerateClusterName()
	p.RegisterCallbacks(id, "update", common.DefaultDB.BroadcastObject)

	// delete k3s node from the selected cluster.
	if err := p.DeleteK3sNode(clusterID, nodeName, instanceId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to delete node : " + clusterID + " node : " + nodeName + "instanceId: " + instanceId + " error: " + err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func CreateCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster"]
	if clusterID == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("clusterID cannot be empty"))
		return
	}

	state, err := common.DefaultDB.GetClusterByID(clusterID)
	if err == nil && state == nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		config := &types.Cluster{}
		if err := json.Unmarshal(body, config); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		provider, err := providers.GetProvider(config.Provider)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to get provider  , error: " + err.Error()))
			return
		}

		err = provider.SetConfig(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		err = provider.MergeClusterOptions()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		id := provider.GenerateClusterName()
		if err = provider.CreateCheck(); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		if err = provider.BindCredential(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// register log callbacks
		provider.RegisterCallbacks(id, "create", common.DefaultDB.BroadcastObject)
		err = provider.CreateK3sCluster()
		if err != nil {
			logrus.Errorf("create cluster error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("cluster %s is existing", clusterID)))
		return
	}
}

func GetClustersList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	clusters, err := cluster.ListClusters(provider)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshal data, error: " + err.Error()))
		return
	}

	cl, err := json.Marshal(clusters)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshal data, error: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cl))
}

func JoinCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster"]
	if clusterID == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("clusterID cannot be empty"))
		return
	}

	state, err := common.DefaultDB.GetClusterByID(clusterID)
	if err != nil || state == nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("cluster %s is not found", clusterID)))
		return
	}

	provider, err := providers.GetProvider(state.Provider)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("provider %s is not found", state.Provider)))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to read join cluster request data ,error:  %s", err.Error())))
		return
	}

	err = provider.SetConfig(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to call SetConfig:  %s", err.Error())))
		return
	}

	err = provider.MergeClusterOptions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to call MergeClusterOptions:  %s", err.Error())))
		return
	}

	id := provider.GenerateClusterName()
	if err = provider.JoinCheck(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to call GenerateClusterName:  %s", err.Error())))
		return
	}

	provider.RegisterCallbacks(id, "update", common.DefaultDB.BroadcastObject)
	err = provider.JoinK3sNode()
	if err != nil {
		logrus.Errorf("join cluster error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to join cluster,error:  %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func ListClusterDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster"]
	kubeCfg := filepath.Join(common.CfgPath, common.KubeCfgFile)

	if clusterID == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("clusterID cannot be empty"))
		return
	}

	state, err := common.DefaultDB.GetClusterByID(clusterID)
	if err != nil || state == nil {
		logrus.Errorf("find cluster error %v", err)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("Not found cluster id %s ", clusterID)))
		return
	}

	provider, err := providers.GetProvider(state.Provider)
	if err != nil {
		logrus.Errorf("failed to get provider %v: %v", state.Provider, err)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to get provider %v: %v", state.Provider, err)))
		return
	}

	provider.SetMetadata(&state.Metadata)
	_ = provider.SetOptions(state.Options)
	isExist, _, err := provider.IsClusterExist()
	if !isExist && err != nil {
		logrus.Errorf("failed to check cluster %s exist, got error: %v ", state.Name, err)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to check cluster %s exist, got error: %v ", state.Name, err)))
		return
	}

	info := provider.DescribeCluster(kubeCfg)
	cl, err := json.Marshal(info)
	if err != nil {
		w.Write([]byte("failed to marshal data, error: " + err.Error()))
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to marshal data, error: %v ", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cl))
}

func PingHost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster"]
	host := vars["host"]

	if clusterID == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("clusterID cannot be empty"))
		return
	}

	if host == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("clusterID cannot be empty"))
		return
	}

	state, err := common.DefaultDB.GetClusterByID(clusterID)
	if err != nil {
		logrus.Errorf("find cluster error %v", err)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("Not found cluster id %s ", clusterID)))
		return
	}

	if state == nil {
		sshInfo := &types.SSH{}
		if err := json.NewDecoder(r.Body).Decode(&sshInfo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))

			return
		}

		hostInfo := strings.Join([]string{host, ":", sshInfo.SSHPort}, "")

		if len(sshInfo.SSHPassword) == 0 && !sshInfo.SSHAgentAuth && len(sshInfo.SSHKeyPath) > 0 {
			sshKey, err := utils.SSHPrivateKeyPath(sshInfo.SSHKeyPath)
			if err != nil {
				return
			}

			sshCert := ""
			if len(sshInfo.SSHCertPath) > 0 {
				result, err := utils.SSHCertificatePath(sshInfo.SSHCertPath)
				if err != nil {
					return
				}
				sshCert = result
			}

			cfg, err := utils.GetSSHConfig(sshInfo.SSHUser, sshKey, sshInfo.SSHKeyPassphrase, sshCert, sshInfo.SSHPassword, 0, sshInfo.SSHAgentAuth)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			// establish connection with SSH server.
			_, err = ssh.Dial("tcp", hostInfo, cfg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))

				return
			}
		} else {
			cfg, err := utils.GetSSHConfig(sshInfo.SSHUser, "", sshInfo.SSHKeyPassphrase, "", sshInfo.SSHPassword, 0, sshInfo.SSHAgentAuth)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))

				return
			}
			// establish connection with SSH server.
			_, err = ssh.Dial("tcp", hostInfo, cfg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}
	} else {
		hostInfo := strings.Join([]string{host, ":", state.SSHPort}, "")

		if len(state.SSHPassword) == 0 && !state.SSHAgentAuth && len(state.SSHKeyPath) > 0 {
			sshKey, err := utils.SSHPrivateKeyPath(state.SSHKeyPath)
			if err != nil {
				return
			}

			sshCert := ""
			if len(state.SSHCertPath) > 0 {
				result, err := utils.SSHCertificatePath(state.SSHCertPath)
				if err != nil {
					return
				}
				sshCert = result
			}

			cfg, err := utils.GetSSHConfig(state.SSHUser, sshKey, state.SSHKeyPassphrase, sshCert, state.SSHPassword, 0, state.SSHAgentAuth)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			// establish connection with SSH server.
			_, err = ssh.Dial("tcp", hostInfo, cfg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))

				return
			}
		} else {
			cfg, err := utils.GetSSHConfig(state.SSHUser, "", state.SSHKeyPassphrase, "", state.SSHPassword, 0, state.SSHAgentAuth)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))

				return
			}
			// establish connection with SSH server.
			_, err = ssh.Dial("tcp", hostInfo, cfg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func updateNodeLabels(contextName string, nodeName string, labelKey string, labelValue string) {
	ctx := context.TODO()
	client, err := cluster.GetClusterConfig(contextName, filepath.Join(common.CfgPath, common.KubeCfgFile))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	labelPatch := fmt.Sprintf(`[{"op":"replace","path":"/metadata/labels/%s","value":"%s" }]`, labelKey, labelValue)
	_, err = client.CoreV1().Nodes().Patch(ctx, nodeName, k8stypes.JSONPatchType, []byte(labelPatch), metav1.PatchOptions{})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
