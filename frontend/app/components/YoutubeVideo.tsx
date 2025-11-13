import { Download, Repeat } from "@mui/icons-material";
import { Button, CircularProgress, SelectChangeEvent, Typography } from "@mui/material";
import { useState } from "react";
import YouTube from "react-youtube";
import YouTubeVideoId from "youtube-video-id";

export interface YoutubeVideoProps {
    videoURL?: string;
}
enum STATUS {
    READY,
    DOWNLOADING,
    FINISHED,
    ERROR
}

export default function YoutubeVideo({ videoURL }: YoutubeVideoProps) {
    const [status, setStatus] = useState<STATUS>(STATUS.READY);

    async function downloadVideo() {
        if (!videoURL) return;
        try {
            setStatus(STATUS.DOWNLOADING);
            const videoID = YouTubeVideoId(videoURL);
            const response = await fetch(`http://localhost:3333/download?id=${videoID}`);
            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const result = await response.text();
            setStatus(STATUS.FINISHED);

        } catch (err) {
            console.error("Download failed:", err);
            setStatus(STATUS.ERROR);
        }
    }

    //Hat hier nichts zu suchen aber bin zu faul es auszulagern
    function getStatusText() {
        switch (status) {
            case STATUS.READY: return <Button variant="contained" startIcon={<Download />} onClick={() => downloadVideo()}>
                Runnerladen
            </Button>;
            case STATUS.DOWNLOADING: return <CircularProgress color="secondary" />;
            case STATUS.FINISHED: return (
                <div>
                    <Typography>Download erfolgreich!</Typography>
                    <Button variant="contained" startIcon={<Download />} onClick={() => downloadVideo()}>Runnerladen</Button>
                </div>
            )
            case STATUS.ERROR: return <Typography>Fehler beim runterladen. Schau mal in die Konsole oder so</Typography>
        }
    }

    const opts = {
        height: "390",
        width: "640",
        playerVars: {
            autoplay: 0,
        },
    };

    if (!videoURL) return null;

    try {
        const videoID = YouTubeVideoId(videoURL);

        if (!videoID) {
            return <Typography color="textPrimary" variant="h3">Bist deppert die ID {videoURL} ist Mist</Typography>;
        }


        return (
            <div className="flex flex-col items-center justify-center w-full gap-4">
                <YouTube videoId={videoID} opts={opts} />
                {getStatusText()}
            </div>
        );

    } catch (e) {
        console.error("Invalid YouTube ID:", e);
        return <Typography color="textPrimary" variant="h3">Bist deppert die ID {videoURL} ist Mist</Typography>;

    }
}
