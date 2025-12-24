import { Download, Repeat } from "@mui/icons-material";
import { Button, CircularProgress, SelectChangeEvent, TextField, Typography } from "@mui/material";
import { useState } from "react";
import YouTube from "react-youtube";
import YouTubeVideoId from "youtube-video-id";

export interface YoutubeVideoProps {
    videoURL?: string;
}


export default function YoutubeVideo({ videoURL }: YoutubeVideoProps) {



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
            return <Typography color="textPrimary" variant="h3">Invalid YouTube URL. Please check and try again.</Typography>;
        }


        return (
            <div className="flex flex-col items-center justify-center w-full gap-4 youtube-preview">
                <YouTube videoId={videoID} opts={opts} />

            </div>
        );

    } catch (e) {
        console.error("Invalid YouTube ID:", e);
        return <Typography color="textPrimary" variant="h3">Invalid YouTube URL. Please check and try again.</Typography>;

    }
}
