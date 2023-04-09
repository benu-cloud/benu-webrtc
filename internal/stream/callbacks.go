package stream

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-webrtc-1.0 gstreamer-sdp-1.0 gstreamer-rtp-1.0
#cgo CFLAGS: -I${SRCDIR}/c
#cgo LDFLAGS: -L${SRCDIR}/c -lstream
#include "stream.h"
*/
import "C"

import (
	"fmt"

	"github.com/benu-cloud/benu-webrtc/pkg/message"
	"github.com/benu-cloud/benu-webrtc/pkg/pkgerrors"
)

//export got_gstreamer_pipeline_error_cb
func got_gstreamer_pipeline_error_cb(from *C.char, message *C.char) {
	if checkStreamInstance() != nil {
		return
	}
	instance.serverGStreamerErrors <- pkgerrors.NewCStreamErrorWithMessage(fmt.Sprintf("error from %s: %s", C.GoString(from), C.GoString(message)))
}

//export got_server_offer_sdp_cb
func got_server_offer_sdp_cb(peerId *C.char, offer *C.char) {
	if checkStreamInstance() != nil {
		return
	}
	fullPid := C.GoString(peerId)
	mtype := fullPid[:1]
	pid := fullPid[1:]
	isVideo := false
	if mtype == "v" {
		isVideo = true
	}
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	// look for peer
	for _, peer := range instance.users {
		if peer.peer_id == pid {
			sdp := C.GoString(offer)
			if isVideo {
				peer.serverSessionDescriptions <- &message.SessionDescriptionPayload{
					From:               peer.peer_id,
					Target:             message.Video,
					SessionDescription: []byte(sdp),
				}
			} else {
				peer.serverSessionDescriptions <- &message.SessionDescriptionPayload{
					From:               peer.peer_id,
					Target:             message.Audio,
					SessionDescription: []byte(sdp),
				}
			}
			return
		}
	}
}

//export got_server_ice_candidate_cb
func got_server_ice_candidate_cb(peerId *C.char, mlineindex C.uint, candidate *C.char) {
	if checkStreamInstance() != nil {
		return
	}
	fullPid := C.GoString(peerId)
	mtype := fullPid[:1]
	pid := fullPid[1:]
	isVideo := false
	if mtype == "v" {
		isVideo = true
	}
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	// look for peer
	for _, peer := range instance.users {
		if peer.peer_id == pid {
			can := C.GoString(candidate)
			if isVideo {
				peer.serverIceCandidates <- &message.IceCandidatePayload{
					From:         peer.peer_id,
					Target:       message.Video,
					IceCandidate: []byte(can),
				}
			} else {
				peer.serverIceCandidates <- &message.IceCandidatePayload{
					From:         peer.peer_id,
					Target:       message.Audio,
					IceCandidate: []byte(can),
				}
			}
			return
		}
	}
}

//export got_client_datachannel_message_cb
func got_client_datachannel_message_cb(peerId *C.char, message *C.char) {
	fmt.Println(4)
}

//export got_webrtc_connection_disconnected_cb
func got_webrtc_connection_disconnected_cb(peerId *C.char) {
	fmt.Println(5)
}
