module github.com/pantheon-systems/container-linux-update-operator

go 1.13

require (
	github.com/blang/semver v3.5.0+incompatible
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/coreos/locksmith v0.6.2-0.20171013225126-ef4279232ecd
	github.com/coreos/pkg v0.0.0-20180108230652-97fdf19511ea
	github.com/godbus/dbus v4.0.0+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/mock v1.2.0
	k8s.io/api v0.0.0-20200131232428-e3a917c59b04
	k8s.io/apimachinery v0.0.0-20200131232151-0cd702f8b7f4
	k8s.io/client-go v0.0.0-20200228043304-076fbc5c36a7
)
