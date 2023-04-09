package message

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    GenericPayload
		wantErr bool
	}{
		{
			name: "valid audio sdp message",
			args: args{
				bytes: []byte(`{
					"type": "sdp",
					"payload": {
						"from": "somepeer",
						"target": "audio",
						"content": {"some sdp": "sdp validity not checked by Go, it is checked by C which is not tested here"}
					}
				}`),
			},
			want: &SessionDescriptionPayload{
				From:               "somepeer",
				Target:             Audio,
				SessionDescription: []byte(`{"some sdp": "sdp validity not checked by Go, it is checked by C which is not tested here"}`),
			},
		},
		{
			name: "valid video sdp message",
			args: args{
				bytes: []byte(`{
					"type": "sdp",
					"payload": {
						"from": "somepeer",
						"target": "video",
						"content": {"some sdp": "sdp validity not checked by Go, it is checked by C which is not tested here"}
					}
				}`),
			},
			want: &SessionDescriptionPayload{
				From:               "somepeer",
				Target:             Video,
				SessionDescription: []byte(`{"some sdp": "sdp validity not checked by Go, it is checked by C which is not tested here"}`),
			},
		},
		{
			name: "valid unknown sdp message",
			args: args{
				bytes: []byte(`{type: "sdp", payload: {
					from: "somepeer",
					target: "something else",
					content: {"some sdp" : "sdp validity not checked by Go, it is checked by C which is not tested here"}
				}}`),
			},
			wantErr: true,
		},
		{
			name: "valid audio ice candidate",
			args: args{
				bytes: []byte(`{
					"type": "ice",
					"payload": {
						"from": "somepeer",
						"target": "audio",
						"content": {"some ice": "ice validity not checked by Go, it is checked by C which is not tested here"}
					}
				}`),
			},
			want: &IceCandidatePayload{
				From:         "somepeer",
				Target:       Audio,
				IceCandidate: []byte(`{"some ice": "ice validity not checked by Go, it is checked by C which is not tested here"}`),
			},
		},
		{
			name: "valid video ice candidate",
			args: args{
				bytes: []byte(`{
					"type": "ice",
					"payload": {
						"from": "somepeer",
						"target": "video",
						"content": {"some ice": "ice validity not checked by Go, it is checked by C which is not tested here"}
					}
				}`),
			},
			want: &IceCandidatePayload{
				From:         "somepeer",
				Target:       Video,
				IceCandidate: []byte(`{"some ice": "ice validity not checked by Go, it is checked by C which is not tested here"}`),
			},
		},
		{
			name: "valid unknown ice candidate",
			args: args{
				bytes: []byte(`{type: "ice", payload: {
					from: "somepeer",
					target: "something else",
					content: {"some ice" : "ice validity not checked by Go, it is checked by C which is not tested here"}
				}}`),
			},
			wantErr: true,
		},
		{
			name: "unknown message type",
			args: args{
				bytes: []byte(`{type: "unknown stuff", payload: {
					from: "somepeer",
					target: "audio",
					content: {"some ice" : "ice validity not checked by Go, it is checked by C which is not tested here"}
				}}`),
			},
			wantErr: true,
		},
		{
			name: "bad sdp payload - no 'from' in payload",
			args: args{
				bytes: []byte(`{type: "sdp", payload: {
					target: "audio",
					content: {"some sdp" : "sdp validity not checked by Go, it is checked by C which is not tested here"}
				}}`),
			},
			wantErr: true,
		},
		{
			name: "bad sdp payload - no 'content' in payload",
			args: args{
				bytes: []byte(`{type: "sdp", payload: {
					from: "somepeer",
					target: "audio",
				}}`),
			},
			wantErr: true,
		},
		{
			name: "bad sdp payload - empty payload",
			args: args{
				bytes: []byte(`{type: "sdp", payload: {}}`),
			},
			wantErr: true,
		},
		{
			name: "bad message with no payload",
			args: args{
				bytes: []byte(`{something: "bad", payload: {}}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	type args struct {
		payload GenericPayload
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "marshal valid audio sdp",
			args: args{
				payload: SessionDescriptionPayload{
					From:               "somepeer",
					Target:             Audio,
					SessionDescription: []byte(`{"some sdp": "sdp validity not checked by Go, it is checked by C which is not tested here"}`),
				},
			},
			want: []byte(`{"type":"sdp","payload":{"from":"somepeer","target":"audio","content":{"some sdp":"sdp validity not checked by Go, it is checked by C which is not tested here"}}}`),
		},
		{
			name: "marshal valid video sdp",
			args: args{
				payload: SessionDescriptionPayload{
					From:               "somepeer",
					Target:             Video,
					SessionDescription: []byte(`{"some sdp": "sdp validity not checked by Go, it is checked by C which is not tested here"}`),
				},
			},
			want: []byte(`{"type":"sdp","payload":{"from":"somepeer","target":"video","content":{"some sdp":"sdp validity not checked by Go, it is checked by C which is not tested here"}}}`),
		},
		{
			name: "marshal valid audio ice candidate",
			args: args{
				payload: IceCandidatePayload{
					From:         "somepeer",
					Target:       Audio,
					IceCandidate: []byte(`{"some ice": "ice validity not checked by Go, it is checked by C which is not tested here"}`),
				},
			},
			want: []byte(`{"type":"ice","payload":{"from":"somepeer","target":"audio","content":{"some ice":"ice validity not checked by Go, it is checked by C which is not tested here"}}}`),
		},
		{
			name: "marshal valid video ice candidate",
			args: args{
				payload: IceCandidatePayload{
					From:         "somepeer",
					Target:       Video,
					IceCandidate: []byte(`{"some ice": "ice validity not checked by Go, it is checked by C which is not tested here"}`),
				},
			},
			want: []byte(`{"type":"ice","payload":{"from":"somepeer","target":"video","content":{"some ice":"ice validity not checked by Go, it is checked by C which is not tested here"}}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
