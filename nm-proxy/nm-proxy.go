package nmproxy

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/gravitl/netclient/nm-proxy/common"
	"github.com/gravitl/netclient/nm-proxy/manager"
	"github.com/gravitl/netclient/nm-proxy/server"
	"github.com/gravitl/netclient/nm-proxy/stun"
)

func Start(ctx context.Context, mgmChan chan *manager.ManagerAction, apiServerAddr string) {
	log.Println("Starting Proxy...")
	common.IsHostNetwork = (os.Getenv("HOST_NETWORK") == "" || os.Getenv("HOST_NETWORK") == "on")
	hInfo := stun.GetHostInfo(apiServerAddr)
	stun.Host = hInfo
	log.Printf("HOSTINFO: %+v", hInfo)
	if IsPublicIP(hInfo.PrivIp) {
		log.Println("Host is public facing!!!")
	}
	// start the netclient proxy server
	err := server.NmProxyServer.CreateProxyServer(0, 0, hInfo.PrivIp.String())
	if err != nil {
		log.Fatal("failed to create proxy: ", err)
	}
	go manager.StartProxyManager(mgmChan)
	server.NmProxyServer.Listen(ctx)

}

// IsPublicIP indicates whether IP is public or not.
func IsPublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate() {
		return false
	}
	return true
}
