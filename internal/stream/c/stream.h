#ifndef STREAM_H
#define STREAM_H
#define GST_USE_UNSTABLE_API
#include <stdbool.h>

typedef enum
{
    SUCCESS,
    ERROR_ENCODER_NOT_SUPPORTED,
    ERROR_PIPELINE_ALREADY_CREATED,
    ERROR_PIPELINE_PARSE_BAD_FORMAT,
    ERROR_PIPELINE_BAD_STATE,
    ERROR_PIPELINE_SET_STATE,
    ERROR_LINKING_PEER,
    ERROR_BAD_PEER_ID,
    ERROR_PIPELINE_DOESNT_EXIST,
    ERROR_BAD_SDP,
} ErrorCode;

typedef enum
{
    STOPPED = 0,
    NONE,
    READY,
    PLAYING,
} PipelineState;

typedef enum
{
    VP9,
    H264,
    NVH264
} VideoEncoder;

typedef struct
{
    // bitrate in bit/s
    unsigned audioBaseBitrate;
    unsigned audioBasePacketLossPct;
    // bitrate in kbit/s
    unsigned videoBaseBitrate;
    unsigned videoBaseFramerate;
    VideoEncoder videoEncoder;
    unsigned videoHeight;
    unsigned videoWidth;
    bool videoShowCursor;
} PipelineOptions;

// callbacks defined in Go
extern void got_gstreamer_pipeline_error_cb(char *from, char *message);
extern void got_server_offer_sdp_cb(char *peerId, char *offer);
extern void got_server_ice_candidate_cb(char *peerId, unsigned int mlineindex, char *candidate);
extern void got_client_datachannel_message_cb(char *peerId, char *message);
extern void got_webrtc_connection_disconnected_cb(char *peerId);

// globally accessible - managed by C
ErrorCode SetupPipeline(PipelineOptions opt);
ErrorCode StartPipeline();
ErrorCode StopPipeline();
ErrorCode AddPeerToPipeline(const char *peer_id);
ErrorCode SetRemoteAnswer(const char *peer_id, const char *answer_sdp);
ErrorCode AddRemoteIceCandidate(const char *peer_id, unsigned int mlineindex, const char *candidate);
ErrorCode RemovePeerFromPipeline(const char *peer_id);

#endif