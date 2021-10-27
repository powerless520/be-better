package model

import DanaData "gitlab.ops.kingnet.com/danadata/go-sdk/com.kingnetdc.danadata.gosdk"

type DanaClient struct {
	Facm 		*DanaData.DanaData
	FacmEvent 	*DanaData.DanaData
}

func (this DanaClient) Close() {
	if this.Facm != nil {
		this.Facm.Close()
	}
	if this.FacmEvent != nil {
		this.FacmEvent.Close()
	}
}