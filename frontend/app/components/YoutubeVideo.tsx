import { Download } from "@mui/icons-material";
import { Button, SelectChangeEvent, Typography } from "@mui/material";
import { useState } from "react";
import YouTube from "react-youtube";
import YouTubeVideoId from "youtube-video-id";

export interface YoutubeVideoProps {
    videoURL?: string;
}


type YoutubeDownloadOptions = {
    quality: string;
}

export default function YoutubeVideo({ videoURL }: YoutubeVideoProps) {
    async function downloadVideo() {
        if (!videoURL) return;
        try {
            const videoID = YouTubeVideoId(videoURL);
            const response = await fetch(`http://localhost:3333/download?id=${videoID}`);
            if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

            const result = await response.text(); // plain string
            console.log("Server says:", result);
            alert(result); // or update state if you want to display it in the UI
        } catch (err) {
            console.error("Download failed:", err);
            alert("Fehler beim Download!");
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
                <Button variant="contained" startIcon={<Download />} onClick={() => downloadVideo()}>
                    Runnerladen
                </Button>
            </div>
        );
    } catch (e) {
        console.error("Invalid YouTube ID:", e);
        return <Typography color="textPrimary" variant="h3">Bist deppert die ID {videoURL} ist Mist</Typography>;

    }
}
