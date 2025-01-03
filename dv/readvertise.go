package dv

import (
	"time"

	enc "github.com/zjkmxy/go-ndn/pkg/encoding"
	"github.com/zjkmxy/go-ndn/pkg/log"
	"github.com/zjkmxy/go-ndn/pkg/ndn"
	mgmt "github.com/zjkmxy/go-ndn/pkg/ndn/mgmt_2022"
	"github.com/zjkmxy/go-ndn/pkg/security"
	"github.com/zjkmxy/go-ndn/pkg/utils"
)

// Received advertisement Interest
func (dv *Router) readvertiseOnInterest(args ndn.InterestHandlerArgs) {
	res := &mgmt.ControlResponse{
		Val: &mgmt.ControlResponseVal{
			StatusCode: 400,
			StatusText: "Failed to execute command",
			Params:     nil,
		},
	}

	defer func() {
		signer := security.NewSha256Signer()
		data, err := dv.engine.Spec().MakeData(
			args.Interest.Name(),
			&ndn.DataConfig{
				ContentType: utils.IdPtr(ndn.ContentTypeBlob),
				Freshness:   utils.IdPtr(1 * time.Second),
			},
			res.Encode(),
			signer)
		if err != nil {
			log.Warnf("readvertise: failed to make response Data: %+v", err)
			return
		}
		args.Reply(data.Wire)
	}()

	// /localhost/nlsr/rib/register/h%0C%07%07%08%05cathyo%01A/params-sha256=a971bb4753691b756cb58239e2585362a154ec6551985133990c8bd2401c466a
	// readvertise:  /localhost/nlsr/rib/unregister/h%0C%07%07%08%05cathyo%01A/params-sha256=026dd595c75032c5101b321fbc11eeb96277661c66bc0564ac7ea1a281ae8210
	iname := args.Interest.Name()
	if len(iname) != 6 {
		log.Warnf("readvertise: invalid interest %s", iname)
		return
	}

	module, cmd, advC := iname[2], iname[3], iname[4]
	if module.String() != "rib" {
		log.Warnf("readvertise: unknown module %s", iname)
		return
	}

	params, err := mgmt.ParseControlParameters(enc.NewBufferReader(advC.Val), false)
	if err != nil || params.Val == nil || params.Val.Name == nil {
		log.Warnf("readvertise: failed to parse advertised name (%s)", err)
		return
	}

	log.Debugf("readvertise: %s %s", cmd, params.Val.Name)
	dv.mutex.Lock()
	defer dv.mutex.Unlock()

	switch cmd.String() {
	case "register":
		dv.pfx.Announce(params.Val.Name)
	case "unregister":
		dv.pfx.Withdraw(params.Val.Name)
	default:
		log.Warnf("readvertise: unknown cmd %s", cmd)
		return
	}

	res.Val.StatusCode = 200
	res.Val.StatusText = "Readvertise command successful"
	res.Val.Params = &mgmt.ControlArgs{
		Name:   params.Val.Name,
		FaceId: utils.IdPtr(uint64(1)), // NFD compatibility
		Origin: utils.IdPtr(uint64(65)),
		Cost:   utils.IdPtr(uint64(0)),
		Flags:  utils.IdPtr(uint64(0)),
	}
}
