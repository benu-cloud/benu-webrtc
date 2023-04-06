/*gstreamer-1.0 gstreamer-webrtc-1.0 gstreamer-sdp-1.0 gstreamer-rtp-1.0*/
#include "stream.h"
#include <gst/webrtc/webrtc.h>
#include <gst/rtp/rtp.h>
#include <glib/gprintf.h>
#include <windows.h>

// === Initialize global variables for file ===
static GstElement *pipeline = NULL;
static PipelineState state = NONE;
static PipelineOptions options;
// === Initialize static functions ===
static void lock();
static void unlock();
static PipelineState getPipelineState();
static bool setPipelineState(PipelineState newState);
static inline int getNumCores();
static void PrintWebRTCStates(GstElement *webrtc);
static void createDotFile();
static ErrorCode createPipeline();
static gboolean on_pipeline_message(GstBus *bus, GstMessage *message, G_GNUC_UNUSED gpointer none);
static void on_connection_state_change(GstElement *webrtc, GParamSpec G_GNUC_UNUSED *pspec, G_GNUC_UNUSED gpointer none);
static void on_negotiation_needed(GstElement *webrtc, G_GNUC_UNUSED gpointer none);
static void on_offer_created(GstPromise *promise, GstElement *webrtc);
static void on_ice_candidate(GstElement G_GNUC_UNUSED *webrtc, guint mlineindex, gchar *candidate, G_GNUC_UNUSED gpointer none);
// used for controls exclusively
static void on_datachannel_message_string(GstWebRTCDataChannel G_GNUC_UNUSED *dc, gchar *msg, G_GNUC_UNUSED gpointer none);
// === Initialize global mutex for all operations ===
static GMutex mutex;
/**
 * @brief lock global mutex
 *
 */
static void lock()
{
    g_mutex_lock(&mutex);
}
/**
 * @brief unlock global mutex
 *
 */
static void unlock()
{
    g_mutex_unlock(&mutex);
}
// === Getters and setters for variables (static functions are private to file) ===
/**
 * @brief Get the Pipeline State object (LOCK MUTEX BEFORE USING THIS)
 *
 * @return PipelineState
 */
static PipelineState getPipelineState()
{
    return state;
}
/**
 * @brief Set the Pipeline State object (LOCK MUTEX BEFORE USING THIS)
 *
 * @param newState
 */
static bool setPipelineState(PipelineState newState)
{
    switch (newState)
    {
    case STOPPED:
        if (gst_element_set_state(pipeline, GST_STATE_NULL) == GST_STATE_CHANGE_FAILURE)
            return false;
        state = newState;
        return true;
    case PLAYING:
        if (gst_element_set_state(pipeline, GST_STATE_PLAYING) == GST_STATE_CHANGE_FAILURE)
            return false;
        state = newState;
        return true;
    case READY:
        if (gst_element_set_state(pipeline, GST_STATE_READY) == GST_STATE_CHANGE_FAILURE)
            return false;
        state = newState;
        return true;
    case NONE: // not supported currently
    default:
        return false;
    }
}
// === Auxiliary functions ===
/**
 * @brief Get the number of cpu cores
 *
 * @return int
 */
static inline int getNumCores()
{
    SYSTEM_INFO sysinfo;
    GetSystemInfo(&sysinfo);
    // TODO: look into issues with >8
    return sysinfo.dwNumberOfProcessors;
}
/**
 * @brief print webrtc connection states
 *
 * @param webrtc
 */
static void PrintWebRTCStates(GstElement *webrtc)
{
    gint icegathstate, iceconstate, sigstate, constate;
    g_object_get(webrtc, "ice-gathering-state", &icegathstate, NULL);
    g_object_get(webrtc, "ice-connection-state", &iceconstate, NULL);
    g_object_get(webrtc, "signaling-state", &sigstate, NULL);
    g_object_get(webrtc, "connection-state", &constate, NULL);
    g_print("ice gath: %d, ice con: %d, signaling: %d, connection: %d\n", icegathstate, iceconstate, sigstate, constate);
}
/**
 * @brief Create a dot file (debug)
 *
 */
static void createDotFile()
{
    GST_DEBUG_BIN_TO_DOT_FILE_WITH_TS(GST_BIN(pipeline), GST_DEBUG_GRAPH_SHOW_ALL, "dotfile");
}
// === Init/deinit functions ===
/**
 * @brief Create global pipeline object
 *
 * @return ErrorCode
 */
static ErrorCode createPipeline()
{
    if (GST_IS_OBJECT(pipeline))
        return ERROR_PIPELINE_ALREADY_CREATED;
    ErrorCode returnVal = SUCCESS;
    char *vcaptureLine;
    char *vencoderLine;
    char *vconverterLine;
    char *acaptureLine;
    char *aencoderLine;
    // capture screen using dx9 (or d3d11 if possible)
    GstElementFactory *d3d11VideoCapture = gst_element_factory_find("d3d11screencapturesrc");
    vcaptureLine = g_strdup_printf(""
                                   "%s do-timestamp=false "
                                   "%s=%s blocksize=16384 ! ",
                                   d3d11VideoCapture == NULL ? "dx9screencapsrc" : "d3d11screencapturesrc",
                                   d3d11VideoCapture == NULL ? "cursor" : "show-cursor",
                                   options.videoShowCursor ? "true" : "false");
    // capture audio with wasapi (or wasapi2 if possible)
    GstElementFactory *wasapi2AudioCapture = gst_element_factory_find("wasapi2src");
    acaptureLine = g_strdup_printf(""
                                   "%s slave-method=none "
                                   "loopback=true low-latency=true provide-clock=false "
                                   "do-timestamp=false ! "
                                   "audio/x-raw,channels=2 ! ",
                                   wasapi2AudioCapture == NULL ? "wasapisrc" : "wasapi2src");
    // encoder parameters
    // TODO: more hardware specific encoding pipelines, VMAF(iqa), more optimization on encoder parameters
    switch (options.videoEncoder)
    {
    case VP9:
        vconverterLine = g_strdup_printf(""
                                         "video/x-raw,framerate=%u/1 ! "
                                         "videoconvert qos=true dither=none n-threads=%d ! "
                                         "video/x-raw,format=I420 ! "
                                         "videoscale qos=true n-threads=%d ! "
                                         "video/x-raw,width=%d,height=%d ! ",
                                         options.videoBaseFramerate,
                                         8 /*getNumCores()*/,
                                         8 /*getNumCores()*/,
                                         options.videoWidth,
                                         options.videoHeight);
        vencoderLine = g_strdup_printf(""
                                       "vp9enc qos=true name=videoencoder buffer-initial-size=500 "
                                       "buffer-optimal-size=600 buffer-size=1500 "
                                       "end-usage=cbr target-bitrate=%u lag-in-frames=0 deadline=1 "
                                       "keyframe-max-dist=%d threads=%d "
                                       "max-intra-bitrate=250 cpu-used=8 static-threshold=1 "
                                       "error-resilient=default row-mt=true "
                                       "! "
                                       "video/x-vp9 ! ",
                                       options.videoBaseBitrate * 1000,
                                       2147483647,
                                       8 /*getNumCores()*/);
        break;
    case H264:
        vconverterLine = g_strdup_printf(""
                                         "video/x-raw,framerate=%u/1 ! "
                                         "videoconvert qos=true dither=none n-threads=%d ! "
                                         "video/x-raw,format=I420 ! "
                                         "videoscale qos=true n-threads=%d ! "
                                         "video/x-raw,width=%d,height=%d ! ",
                                         options.videoBaseFramerate,
                                         8 /*getNumCores()*/,
                                         8 /*getNumCores()*/,
                                         options.videoWidth,
                                         options.videoHeight);
        // TODO: look into high-444 in case of moving away from browser
        // profile-level-id from  https://www.iana.org/assignments/media-types/video/H264-SVC
        vencoderLine = g_strdup_printf(""
                                       "x264enc qos=true name=videoencoder vbv-buf-capacity=750 "
                                       "bitrate=%u sliced-threads=true byte-stream=false "
                                       "speed-preset=veryfast key-int-max=%d threads=%d "
                                       "tune=zerolatency b-adapt=false ref=1 psy-tune=ssim bframes=0 "
                                       "! "
                                       "video/x-h264,profile=high,stream-format=avc ! ",
                                       options.videoBaseBitrate,
                                       0,
                                       8 /*getNumCores()*/);
        break;
    case NVH264:
        vconverterLine = g_strdup_printf(""
                                         "video/x-raw(memory:D3D11Memory),framerate=%u/1 ! "
                                         "d3d11convert qos=true ! "
                                         "video/x-raw(memory:D3D11Memory),format=I420 ! "
                                         "d3d11scale qos=true ! "
                                         "video/x-raw(memory:D3D11Memory),width=%d,height=%d ! "
                                         "d3d11download qos=true ! ",
                                         options.videoBaseFramerate,
                                         options.videoWidth,
                                         options.videoHeight);
        vencoderLine = g_strdup_printf(""
                                       "nvh264enc qos=true name=videoencoder bitrate=%u "
                                       "vbv-buffer-size=1300 bframes=0 b-adapt=false rc-lookahead=0 "
                                       "zerolatency=true preset=low-latency-hq rc-mode=cbr "
                                       "gop-size=%d "
                                       "! "
                                       "video/x-h264,profile=high ! ",
                                       options.videoBaseBitrate,
                                       -1);
        break;
    default:
        returnVal = ERROR_ENCODER_NOT_SUPPORTED;
        goto done;
    }
    aencoderLine = g_strdup_printf(""
                                   "audioconvert dithering=none ! "
                                   "opusenc name=audioencoder bitrate=%u hard-resync=true "
                                   "bandwidth=fullband audio-type=restricted-lowdelay "
                                   "inband-fec=true packet-loss-percentage=%u ! ",
                                   options.audioBaseBitrate,
                                   options.audioBasePacketLossPct);
    GError *error = NULL;
    char *basePipelineString;
    basePipelineString = g_strdup_printf(""
                                         // capture video and create rtp packets
                                         // capture video
                                         "%s"
                                         // debug (time overlay)
                                         // "timeoverlay !"
                                         // set framerate, color format and scale
                                         "%s"
                                         // encode
                                         "%s"
                                         // tee for sending rtp packets to potentially many rtc clients
                                         "tee name=videoenctee "
                                         // a copy to fakesink for prerolling early (might not be needed)
                                         "videoenctee. ! "
                                         "queue flush-on-eos=true leaky=downstream silent=true ! "
                                         "fakesink "
                                         // capture audio and create rtp packets
                                         // capture audio
                                         "%s"
                                         // encode
                                         "%s"
                                         // tee for sending rtp packets to potentially many rtc clients
                                         "tee name=audioenctee "
                                         // a copy to fakesink for prerolling early (might not be needed)
                                         "audioenctee. ! "
                                         "queue flush-on-eos=true leaky=downstream silent=true ! "
                                         "fakesink ",
                                         vcaptureLine,
                                         vconverterLine,
                                         vencoderLine,
                                         acaptureLine,
                                         aencoderLine);
    g_free(vencoderLine);
    g_free(vconverterLine);
    g_free(aencoderLine);
    pipeline = gst_parse_launch(basePipelineString, &error);
    // take ownership of floating ref
    if (G_IS_INITIALLY_UNOWNED (pipeline))
        g_object_ref_sink (pipeline);
    g_free(basePipelineString);
    if (error)
    {
        // gst_println(error->message);
        g_error_free(error);
        returnVal = ERROR_PIPELINE_PARSE_BAD_FORMAT;
        goto done;
    }
done:
    g_free(vcaptureLine);
    g_free(acaptureLine);
    return returnVal;
}
/**
 * @brief set up pipeline and put it in the READY state. MUST be run only once and before any other functions
 *
 * @param opt
 * @return ErrorCode
 */
ErrorCode SetupPipeline(PipelineOptions opt)
{
    if (getPipelineState() != NONE)
    {
        return ERROR_PIPELINE_BAD_STATE;
    }
    // init gstreamer
    gst_init(NULL, NULL);
    ErrorCode returnVal = SUCCESS;
    // lock mutex
    lock();

    // set global options
    options = opt;

    // create pipeline
    returnVal = createPipeline();
    if (returnVal != SUCCESS)
        goto done;

    // set pipeline state to READY
    if (!setPipelineState(READY))
    {
        returnVal = ERROR_PIPELINE_SET_STATE;
        goto done;
    }

    // get bus
    GstBus *bus;
    bus = gst_pipeline_get_bus(GST_PIPELINE(pipeline));
    gst_bus_add_watch(bus, (GstBusFunc)on_pipeline_message, NULL);
    gst_object_unref(bus);

done:
    unlock();
    return returnVal;
}
/**
 * @brief start the pipeline and put it in the PLAYING state
 * 
 * @return ErrorCode 
 */
ErrorCode StartPipeline()
{
    ErrorCode returnVal = SUCCESS;
    lock();
    // check if state is valid
    switch (getPipelineState())
    {
    case NONE:
        returnVal = ERROR_PIPELINE_DOESNT_EXIST;
        goto done;
    case STOPPED:
    case PLAYING:
        returnVal = ERROR_PIPELINE_BAD_STATE;
        goto done;
    case READY:
        break;
    }

    // set pipeline state to PLAYING
    if (!setPipelineState(PLAYING))
    {
        returnVal = ERROR_PIPELINE_SET_STATE;
        goto done;
    }
done:
    unlock();
    return returnVal;
}
/**
 * @brief irreversibly stop the pipeline and put it in the STOPPED state, as well as freeing the pipeline object
 * remember to remove all peers from the pipeline before calling this function
 * 
 * @return ErrorCode 
 */
ErrorCode StopPipeline()
{
    ErrorCode returnVal = SUCCESS;
    lock();
    // check if state is valid
    switch (getPipelineState())
    {
    case NONE:
        returnVal = ERROR_PIPELINE_DOESNT_EXIST;
        goto done;
    case STOPPED:
        returnVal =  ERROR_PIPELINE_BAD_STATE;
        goto done;
    case PLAYING:
    case READY:
        break;
    }
    // TODO: send eos to pipeline
    if (!setPipelineState(STOPPED))
    {
        returnVal = ERROR_PIPELINE_SET_STATE;
        goto done;
    }
    /* Free resources */
    gst_object_unref(pipeline);
    return SUCCESS;
done:
    unlock();
    return returnVal;
}

// === Core functions ===
/**
 * @brief add a webrtc peer to the pipeline
 * 
 * @param peer_id 
 * @return ErrorCode 
 */
ErrorCode AddPeerToPipeline(const char *peer_id)
{
    ErrorCode returnVal = SUCCESS;
    lock();
    // peer id names are used a lot
    // so define them now and free in goto
    char *peer_id_vname = g_strdup_printf("v%s", peer_id);
    char *peer_id_aname = g_strdup_printf("a%s", peer_id);

    // check if state is valid
    switch (getPipelineState())
    {
    case NONE:
        returnVal = ERROR_PIPELINE_DOESNT_EXIST;
        goto done;
    case STOPPED:
    case READY:
        returnVal = ERROR_PIPELINE_BAD_STATE;
        goto done;
    case PLAYING:
        break;
    }

    // set parse / payloader settings
    char *vParserPayloader;
    char *aPayloader;
    switch (options.videoEncoder)
    {
    case VP9:
        vParserPayloader = g_strdup_printf(""
                                           "vp9parse ! "
                                           "rtpvp9pay name=videopay picture-id-mode=15-bit ! "
                                           "application/x-rtp,clock-rate=90000,media=video,encoding-name=VP9,payload=%d",
                                           123);
        break;
    case H264:
    case NVH264:
        vParserPayloader = g_strdup_printf(""
                                           "h264parse ! "
                                           "rtph264pay name=videopay config-interval=-1 aggregate-mode=zero-latency ! "
                                           "application/x-rtp,clock-rate=90000,media=video,encoding-name=H264,payload=%d",
                                           123);
        break;
    default:
        returnVal = ERROR_ENCODER_NOT_SUPPORTED;
        goto done;
    }

    aPayloader = g_strdup_printf(""
                                 "rtpopuspay name=audiopay ! "
                                 "application/x-rtp,clock-rate=48000,media=audio,encoding-name=OPUS,payload=%d,"
                                 "stereo=(string)1,minptime=(string)10,rtx-time=(string)125,useinbandfec=(string)1",
                                 97);

    // create webrtc pipelines
    char *vwebrtcLine;
    char *awebrtcLine;
    vwebrtcLine = g_strdup_printf(""
                                  "queue leaky=downstream silent=true max-size-buffers=0 "
                                  "max-size-bytes=0 max-size-time=1000000000 flush-on-eos=true ! "
                                  "%s ! "
                                  "webrtcbin name=webrtc stun-server=stun://stun.l.google.com:19302 "
                                  "bundle-policy=max-compat latency=1",
                                  vParserPayloader);

    awebrtcLine = g_strdup_printf(""
                                  "queue leaky=downstream silent=true max-size-buffers=0 "
                                  "max-size-bytes=0 max-size-time=1000000000 flush-on-eos=true ! "
                                  "%s ! "
                                  "webrtcbin name=webrtc stun-server=stun://stun.l.google.com:19302 "
                                  "bundle-policy=max-compat latency=1",
                                  aPayloader);
    g_free(vParserPayloader);
    g_free(aPayloader);

    // create webrtc wrapper bins (floating refs are taken ownership in gst_bin_add)
    GstElement *videoWebrtcbin, *audioWebrtcbin;
    videoWebrtcbin = gst_parse_bin_from_description(vwebrtcLine, TRUE, NULL);
    g_free(vwebrtcLine);
    if (videoWebrtcbin == NULL)
    {
        returnVal = ERROR_PIPELINE_PARSE_BAD_FORMAT;
        goto done;
    }
    audioWebrtcbin = gst_parse_bin_from_description(awebrtcLine, TRUE, NULL);
    g_free(awebrtcLine);
    if (audioWebrtcbin == NULL)
    {
        returnVal = ERROR_PIPELINE_PARSE_BAD_FORMAT;
        goto done;
    }
    
    // set element names
    gst_element_set_name(videoWebrtcbin, peer_id_vname);
    gst_element_set_name(audioWebrtcbin, peer_id_aname);

    // enable all extensions manually
    // adding according to chrome support
    // updated for 1.21.x prerelease
    int i;
    GList *exts_list, *ext;
    // video
    GstElement *videopay = gst_bin_get_by_name(GST_BIN(videoWebrtcbin), "videopay");
    g_warn_if_fail(videopay != NULL);
    exts_list = gst_rtp_get_header_extension_list();
    for (ext = exts_list, i = 1; ext; ext = ext->next)
    {
        GstElementFactory *extension_factory = ext->data;
        GstRTPHeaderExtension *extension = GST_RTP_HEADER_EXTENSION_CAST(gst_element_factory_create(extension_factory, NULL));
        const char *uri = gst_rtp_header_extension_get_uri(extension);
        // video-only
        if (g_str_match_string("color-space", uri, FALSE) ||
            g_str_match_string("rtp-stream-id", uri, FALSE) ||
            g_str_match_string("sdes:mid", uri, FALSE) ||
            g_str_match_string("transport-wide-cc", uri, FALSE))
        {
            gst_rtp_header_extension_set_id(extension, i);
            i++;
            g_signal_emit_by_name(videopay, "add-extension", extension);
        }
        else
        {
            g_info("uri %s not added to video payloader", uri);
        }
    }
    gst_plugin_feature_list_free(exts_list);
    gst_object_unref(videopay);
    // audio
    GstElement *audiopay = gst_bin_get_by_name(GST_BIN(audioWebrtcbin), "audiopay");
    g_warn_if_fail(audiopay != NULL);
    exts_list = gst_rtp_get_header_extension_list();
    for (ext = exts_list, i = 1; ext; ext = ext->next)
    {
        GstElementFactory *extension_factory = ext->data;
        GstRTPHeaderExtension *extension = GST_RTP_HEADER_EXTENSION_CAST(gst_element_factory_create(extension_factory, NULL));
        const char *uri = gst_rtp_header_extension_get_uri(extension);
        // audio-only
        if (g_str_match_string("audio-level", uri, FALSE) ||
            g_str_match_string("sdes:mid", uri, FALSE) ||
            g_str_match_string("transport-wide-cc", uri, FALSE))
        {
            gst_rtp_header_extension_set_id(extension, i);
            i++;
            g_signal_emit_by_name(audiopay, "add-extension", extension);
        }
        else
        {
            g_info("uri %s not added to audio payloader", uri);
        }
    }
    gst_plugin_feature_list_free(exts_list);
    gst_object_unref(audiopay);
    // add to pipeline - ownership is transferred to parent
    g_warn_if_fail(gst_bin_add(GST_BIN(pipeline), videoWebrtcbin));
    g_warn_if_fail(gst_bin_add(GST_BIN(pipeline), audioWebrtcbin));
    // link to main encoders
    GstPad *srcpad;
    GstElement *tee;
    int ret;

    // link vid
    tee = gst_bin_get_by_name(GST_BIN(pipeline), "videoenctee");
    g_assert_nonnull(tee);
    srcpad = gst_element_request_pad_simple(tee, "src_%u");
    g_assert_nonnull(srcpad);
    ret = gst_pad_link(srcpad, videoWebrtcbin->sinkpads->data);
    gst_object_unref(tee);
    gst_object_unref(srcpad);
    if (ret != GST_PAD_LINK_OK)
    {
        returnVal = ERROR_LINKING_PEER;
        goto done;
    }
    // link aud
    tee = gst_bin_get_by_name(GST_BIN(pipeline), "audioenctee");
    g_assert_nonnull(tee);
    srcpad = gst_element_request_pad_simple(tee, "src_%u");
    g_assert_nonnull(srcpad);
    ret = gst_pad_link(srcpad, audioWebrtcbin->sinkpads->data);
    gst_object_unref(tee);
    gst_object_unref(srcpad);
    if (ret != GST_PAD_LINK_OK)
    {
        returnVal = ERROR_LINKING_PEER;
        goto done;
    }
    // extract webrtcbin from wrapper
    GstElement *vwebrtcbin = gst_bin_get_by_name(GST_BIN(videoWebrtcbin), "webrtc");
    g_assert_nonnull(vwebrtcbin);
    GstElement *awebrtcbin = gst_bin_get_by_name(GST_BIN(audioWebrtcbin), "webrtc");
    g_assert_nonnull(awebrtcbin);
    // set transceiver settings
    GArray *transceivers;
    GstWebRTCRTPSender *sender;
    GstWebRTCRTPTransceiver *trans;
    // video
    g_signal_emit_by_name(vwebrtcbin, "get-transceivers", &transceivers);
    g_assert(transceivers != NULL && transceivers->len > 0);

    trans = g_array_index(transceivers, GstWebRTCRTPTransceiver *, 0);
    g_object_set(trans, "direction", GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_SENDONLY, NULL);
    g_object_set(trans, "do-nack", TRUE, NULL);
    g_object_get(trans, "sender", &sender, NULL);
    gst_webrtc_rtp_sender_set_priority(sender, GST_WEBRTC_PRIORITY_TYPE_HIGH);
    g_object_unref(sender);
    g_object_unref(trans);
    g_array_unref(transceivers);

    // audio
    g_signal_emit_by_name(awebrtcbin, "get-transceivers", &transceivers);
    g_assert(transceivers != NULL && transceivers->len > 0);

    trans = g_array_index(transceivers, GstWebRTCRTPTransceiver *, 0);
    g_object_set(trans, "direction", GST_WEBRTC_RTP_TRANSCEIVER_DIRECTION_SENDONLY, NULL);
    g_object_get(trans, "sender", &sender, NULL);
    g_object_set(trans, "fec-type", GST_WEBRTC_FEC_TYPE_ULP_RED, NULL);
    g_object_set(trans, "fec-percentage", options.audioBasePacketLossPct, NULL);
    gst_webrtc_rtp_sender_set_priority(sender, GST_WEBRTC_PRIORITY_TYPE_HIGH);
    g_object_unref(sender);
    g_object_unref(trans);
    g_array_unref(transceivers);
    // add signal handlers
    g_signal_connect(vwebrtcbin, "notify::connection-state", G_CALLBACK(on_connection_state_change), NULL);
    g_signal_connect(awebrtcbin, "notify::connection-state", G_CALLBACK(on_connection_state_change), NULL);
    g_signal_connect(vwebrtcbin, "on-negotiation-needed", G_CALLBACK(on_negotiation_needed), NULL);
    g_signal_connect(awebrtcbin, "on-negotiation-needed", G_CALLBACK(on_negotiation_needed), NULL);
    g_signal_connect(vwebrtcbin, "on-ice-candidate", G_CALLBACK(on_ice_candidate), NULL);
    g_signal_connect(awebrtcbin, "on-ice-candidate", G_CALLBACK(on_ice_candidate), NULL); 
    
    // sync states with parent
    ret = gst_element_sync_state_with_parent(audioWebrtcbin);
    g_assert_true(ret);
    ret = gst_element_sync_state_with_parent(videoWebrtcbin);
    g_assert_true(ret);
    // datachannel is added to audio webrtcbin after state set to playing
    GstWebRTCDataChannel *datachannel;
    GstStructure *datachannelSettings;
    datachannelSettings = gst_structure_new("settings",
                                            "ordered", G_TYPE_BOOLEAN, TRUE,
                                            "priority", GST_TYPE_WEBRTC_PRIORITY_TYPE, GST_WEBRTC_PRIORITY_TYPE_HIGH,
                                            "max-retransmits", G_TYPE_INT, 1,
                                            NULL);
    g_signal_emit_by_name(awebrtcbin, "create-data-channel", "controls", datachannelSettings, &datachannel);
    g_signal_connect(datachannel, "on-message-string", G_CALLBACK(on_datachannel_message_string), NULL);
    // add a reference since the next function doesn't take ownership
    g_object_ref(datachannel);
    // store the datachannel in the webrtcbin
    g_object_set_qdata_full(G_OBJECT(awebrtcbin), g_quark_from_static_string("datachannel-controls"), datachannel, g_object_unref);

    gst_structure_free(datachannelSettings);
    g_object_unref(datachannel);
    gst_object_unref(vwebrtcbin);
    gst_object_unref(awebrtcbin);
    gst_object_unref(audioWebrtcbin);
    gst_object_unref(videoWebrtcbin);
done:
    unlock();
    g_free(peer_id_aname);
    g_free(peer_id_vname);
    return returnVal;
}
/**
 * @brief Set the remove answer for peer
 * Make sure peer exists when using this
 * 
 * @param peer_id
 * @param answer_sdp
 * @return ErrorCode
 */
ErrorCode SetRemoteAnswer(const char *peer_id, const char *answer_sdp)
{
    ErrorCode returnVal = SUCCESS;
    lock();

    // check if state is valid
    switch (getPipelineState())
    {
    case NONE:
        returnVal = ERROR_PIPELINE_DOESNT_EXIST;
        goto done;
    case STOPPED:
    case READY:
        returnVal = ERROR_PIPELINE_BAD_STATE;
        goto done;
    case PLAYING:
        break;
    }

    GstSDPMessage *sdp;
    GstPromise *promise;
    GstWebRTCSessionDescription *answer;
    GstElement *webrtcbin, *webrtc;
    int ret;

    // get peer
    webrtc = gst_bin_get_by_name(GST_BIN(pipeline), peer_id);
    if (!GST_IS_ELEMENT(webrtc))
    {
        returnVal = ERROR_BAD_PEER_ID;
        goto done;
    }

    // create an sdp message and parse into it
    ret = gst_sdp_message_new(&sdp);
    g_assert_cmphex(ret, ==, GST_SDP_OK);
    ret = gst_sdp_message_parse_buffer((guint8 *)answer_sdp, strlen(answer_sdp), sdp);
    if (ret != GST_SDP_OK)
    {
        returnVal = ERROR_BAD_SDP;
        goto done;
    }
    // sdp ownership is transferred
    answer = gst_webrtc_session_description_new(GST_WEBRTC_SDP_TYPE_ANSWER, sdp);
    g_warn_if_fail(answer != NULL);

    webrtcbin = gst_bin_get_by_name(GST_BIN(webrtc), "webrtc");

    promise = gst_promise_new();
    g_signal_emit_by_name(webrtcbin, "set-remote-description", answer, promise);
    gst_object_unref(webrtcbin);
    gst_object_unref(webrtc);
    gst_promise_interrupt(promise);
    gst_promise_unref(promise);
    gst_webrtc_session_description_free(answer);
done:
    unlock();
    return returnVal;
}
/**
 * @brief Set the Remote ice candidate for peer
 * Make sure peer exists when using this
 *
 * @param peer_id
 * @param mlineindex
 * @param candidate
 * @return ErrorCode
 */
ErrorCode AddRemoteIceCandidate(const char *peer_id, unsigned int mlineindex, const char *candidate)
{
    ErrorCode returnVal = SUCCESS;
    lock();

    // check if state is valid
    switch (getPipelineState())
    {
    case NONE:
        returnVal = ERROR_PIPELINE_DOESNT_EXIST;
        goto done;
    case STOPPED:
    case READY:
        returnVal = ERROR_PIPELINE_BAD_STATE;
        goto done;
    case PLAYING:
        break;
    }

    GstElement *webrtc, *webrtcbin;

    webrtc = gst_bin_get_by_name(GST_BIN(pipeline), peer_id);
    if (!GST_IS_ELEMENT(webrtc))
    {
        returnVal = ERROR_BAD_PEER_ID;
        goto done;
    }
    webrtcbin = gst_bin_get_by_name(GST_BIN(webrtc), "webrtc");
    g_signal_emit_by_name(webrtcbin, "add-ice-candidate", mlineindex, candidate);

    gst_object_unref(webrtcbin);
    gst_object_unref(webrtc);
done:
    unlock();
    return returnVal;
}
/**
 * @brief Remove a peer from pipeline and free its resources
 * Make sure peer exists when using this
 *
 * @param peer_id
 * @return ErrorCode
 */
ErrorCode RemovePeerFromPipeline(const char *peer_id)
{
    ErrorCode returnVal = SUCCESS;
    lock();

    char *peer_id_vname;
    char *peer_id_aname;
    peer_id_vname = g_strdup_printf("v%s", peer_id);
    peer_id_aname = g_strdup_printf("a%s", peer_id);

    // check if state is valid
    switch (getPipelineState())
    {
    case NONE:
        returnVal = ERROR_PIPELINE_DOESNT_EXIST;
        goto done;
    case STOPPED:
    case READY:
        returnVal = ERROR_PIPELINE_BAD_STATE;
        goto done;
    case PLAYING:
        break;
    }
    
    GstElement *audioWebrtcbin, *videoWebrtcbin;

    videoWebrtcbin = gst_bin_get_by_name(GST_BIN(pipeline), peer_id_vname);
    if (!GST_IS_ELEMENT(videoWebrtcbin))
    {
        returnVal = ERROR_BAD_PEER_ID;
        goto done;
    }
    audioWebrtcbin = gst_bin_get_by_name(GST_BIN(pipeline), peer_id_aname);
    if (!GST_IS_ELEMENT(audioWebrtcbin))
    {
        returnVal = ERROR_BAD_PEER_ID;
        goto done;
    }

    GstElement *tee;
    GstPad *teepad;
    /* tear down video branch */
    tee = gst_bin_get_by_name(GST_BIN(pipeline), "videoenctee");
    g_assert_nonnull(tee);
    teepad = gst_pad_get_peer(videoWebrtcbin->sinkpads->data);
    gst_element_release_request_pad(tee, teepad);
    gst_object_unref(teepad);
    gst_object_unref(tee);

    /* tear down audio branch */
    tee = gst_bin_get_by_name(GST_BIN(pipeline), "audioenctee");
    g_assert_nonnull(tee);
    teepad = gst_pad_get_peer(audioWebrtcbin->sinkpads->data);
    gst_element_release_request_pad(tee, teepad);
    gst_object_unref(teepad);
    gst_object_unref(tee);


    // remove webrtcbin from webrtc wrappers
    GstElement *webrtc;

    // video
    webrtc = gst_bin_get_by_name(GST_BIN(videoWebrtcbin), "webrtc");
    g_assert_nonnull(webrtc);
    // remove signal handlers
    g_signal_handlers_disconnect_by_func(webrtc, G_CALLBACK(on_connection_state_change), NULL);
    g_signal_handlers_disconnect_by_func(webrtc, G_CALLBACK(on_negotiation_needed), NULL);
    g_signal_handlers_disconnect_by_func(webrtc, G_CALLBACK(on_ice_candidate), NULL);
    // remove bin
    g_warn_if_fail(gst_element_set_state(webrtc, GST_STATE_NULL));
    // also unrefs
    g_warn_if_fail(gst_bin_remove(GST_BIN(videoWebrtcbin), webrtc));
    // gst_object_unref(webrtc);

    // audio
    webrtc = gst_bin_get_by_name(GST_BIN(audioWebrtcbin), "webrtc");
    g_assert_nonnull(webrtc);
    // remove signal handlers
    g_signal_handlers_disconnect_by_func(webrtc, G_CALLBACK(on_connection_state_change), NULL);
    g_signal_handlers_disconnect_by_func(webrtc, G_CALLBACK(on_negotiation_needed), NULL);
    g_signal_handlers_disconnect_by_func(webrtc, G_CALLBACK(on_ice_candidate), NULL);
     // stop datachannel(s)
    GstWebRTCDataChannel* datachannel;
    datachannel = GST_WEBRTC_DATA_CHANNEL(g_object_get_qdata(G_OBJECT(webrtc), g_quark_from_static_string("datachannel-controls")));
    g_assert_nonnull(datachannel);
    g_signal_handlers_disconnect_by_func(datachannel, G_CALLBACK(on_datachannel_message_string), NULL);
    gst_webrtc_data_channel_close(datachannel);
    gst_object_unref(datachannel);
    // free the datachannel object, since set_qdata_full was used freeing is done automatically
    g_object_set_qdata(G_OBJECT(webrtc), g_quark_from_static_string("datachannel-controls"), NULL);
    // remove bin
    g_warn_if_fail(gst_element_set_state(webrtc, GST_STATE_NULL));
    // also unrefs
    g_warn_if_fail(gst_bin_remove(GST_BIN(audioWebrtcbin), webrtc));
    // gst_object_unref(webrtc);

    // set states to null
    g_warn_if_fail(gst_element_set_state(videoWebrtcbin, GST_STATE_NULL));
    g_warn_if_fail(gst_element_set_state(audioWebrtcbin, GST_STATE_NULL));

    // remove elements from pipeline, also unrefs
    g_warn_if_fail(gst_bin_remove(GST_BIN(pipeline), videoWebrtcbin));
    g_warn_if_fail(gst_bin_remove(GST_BIN(pipeline), audioWebrtcbin));
    // gst_object_unref(audioWebrtcbin);
    // gst_object_unref(videoWebrtcbin);

done:
    g_free(peer_id_vname);
    g_free(peer_id_aname);
    unlock();
    return returnVal;
}
// === Callbacks and event handlers ===
/**
 * @brief callback for messages on the pipeline bus
 * 
 * @param bus 
 * @param message 
 * @param none 
 * @return gboolean 
 */
static gboolean on_pipeline_message(GstBus *bus, GstMessage *message, G_GNUC_UNUSED gpointer none)
{
    g_printerr("got message");
    switch (GST_MESSAGE_TYPE(message))
    {
    case GST_MESSAGE_EOS:
        gst_println("got eof");
        // stop the pipeline
        g_warn_if_fail(StopPipeline());
        gst_object_unref(bus);
        // remove bus watch
        return FALSE;
    case GST_MESSAGE_ERROR:
    {
        GError *err;
        gchar *debug;
        gst_message_parse_error(message, &err, &debug);
        got_gstreamer_pipeline_error_cb(GST_OBJECT_NAME(message->src), err->message);
        g_error_free(err);
        g_free(debug);
        // stop the pipeline
        g_warn_if_fail(StopPipeline());
        gst_object_unref(bus);
        break;
    }
    default:
        // unhandled message
        break;
    }
    return TRUE;
}
/**
 * @brief callback for connection state changes
 * 
 * @param webrtc 
 * @param pspec 
 * @param none 
 */
static void on_connection_state_change(GstElement *webrtc, GParamSpec G_GNUC_UNUSED *pspec, G_GNUC_UNUSED gpointer none)
{
    lock();
    switch (getPipelineState())
    {
    case NONE:
    case STOPPED:
    case READY:
        goto done;
    case PLAYING:
        break;
    }
    gint constate;
    g_object_get(webrtc, "connection-state", &constate, NULL);
    if (constate == GST_WEBRTC_PEER_CONNECTION_STATE_DISCONNECTED)
    {
        // get name of peer
        GstElement *parent = GST_ELEMENT(gst_element_get_parent(webrtc));
        gchar *peerId = gst_element_get_name(parent);
        got_webrtc_connection_disconnected_cb(peerId);
        g_free(peerId);
        gst_object_unref(parent);
    }
done:
    unlock();
}
/**
 * @brief callback to notify that negotiation required
 * 
 * @param webrtc 
 * @param none 
 */
static void on_negotiation_needed(GstElement *webrtc, G_GNUC_UNUSED gpointer none)
{
    lock();
    switch (getPipelineState())
    {
    case NONE:
    case STOPPED:
    case READY:
        goto done;
    case PLAYING:
        break;
    }

    GstPromise *promise;
    promise = gst_promise_new_with_change_func((GstPromiseChangeFunc)on_offer_created, (gpointer)webrtc, NULL);
    g_signal_emit_by_name(webrtc, "create-offer", NULL, promise);
done:
    unlock();
}
/**
 * @brief callback to create an offer and send it to Go
 * 
 * @param promise 
 * @param webrtc 
 */
static void on_offer_created(GstPromise *promise, GstElement *webrtc)
{
    lock();
    char *sdp_string;
    GstWebRTCSessionDescription *offer;
    // GstSDPMedia *offer_media;
    const GstStructure *reply;
    if (gst_promise_wait(promise) != GST_PROMISE_RESULT_REPLIED)
    {
        g_warning("on_offer_created promise didn't reply");
        goto done;
    }
    reply = gst_promise_get_reply(promise);
    // no need  to free offer, its a reference
    gst_structure_get(reply, "offer", GST_TYPE_WEBRTC_SESSION_DESCRIPTION, &offer, NULL);
    // free promise since we have ownership
    gst_promise_unref(promise);

    // modify offer
    // offer_media = (GstSDPMedia*) gst_sdp_message_get_media(offer->sdp, 0);
    // gst_sdp_media_add_attribute(offer_media, "fmtp", "stereo=1");
    // gst_sdp_media_free(offer_media);
    promise = gst_promise_new();
    g_signal_emit_by_name(webrtc, "set-local-description", offer, promise);
    gst_promise_interrupt(promise);
    gst_promise_unref(promise);

    // send to external func
    sdp_string = gst_sdp_message_as_text(offer->sdp);

    // get name of peer
    GstElement *parent = GST_ELEMENT(gst_element_get_parent(webrtc));
    gchar *peerId = gst_element_get_name(parent);
    got_server_offer_sdp_cb(peerId, sdp_string);
    g_free(peerId);
    gst_object_unref(parent);

    g_free(sdp_string);
done:
    unlock();
}
/**
 * @brief callback to receive ice candidate and send it to peer
 * 
 * @param webrtc 
 * @param mlineindex 
 * @param candidate 
 * @param none 
 */
static void on_ice_candidate(GstElement G_GNUC_UNUSED *webrtc, guint mlineindex, gchar *candidate, G_GNUC_UNUSED gpointer none)
{
    lock();
    switch (getPipelineState())
    {
    case NONE:
    case STOPPED:
    case READY:
        goto done;
    case PLAYING:
        break;
    }
    // get name of peer
    GstElement *parent = GST_ELEMENT(gst_element_get_parent(webrtc));
    gchar *peerId = gst_element_get_name(parent);
    got_server_ice_candidate_cb(peerId, mlineindex, candidate);
    g_free(peerId);
    gst_object_unref(parent);
done:
    unlock();
}
/**
 * @brief callback to receive user datachannel messages and send it to Go
 * 
 * @param dc 
 * @param msg 
 * @param none 
 */
static void on_datachannel_message_string(GstWebRTCDataChannel G_GNUC_UNUSED *dc, gchar *msg, G_GNUC_UNUSED gpointer none)
{
    // ! No need to check state, since datachannel control messages need to be processed as soon as possible
    // switch (getPipelineState())
    // {
    // case NONE: case STOPPED: case READY:
    //     return;
    // }
    // get name of peer
    GstElement *parent = GST_ELEMENT(g_object_get_qdata(G_OBJECT(dc), g_quark_from_static_string("datachannel-controls")));
    gchar *peerId = gst_element_get_name(parent);
    got_client_datachannel_message_cb(peerId, msg);
    g_free(peerId);
}