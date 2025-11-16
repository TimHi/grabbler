'use client'

import YoutubeVideo from "./components/YoutubeVideo";
import { useState } from "react";
import VideoForm from "./components/VideoForm";
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme';
import { Button, CircularProgress, Typography } from "@mui/material";
import YouTubeVideoId from "youtube-video-id";
import { Download } from "@mui/icons-material";

enum STATUS {
  READY,
  DOWNLOADING,
  FINISHED,
  ERROR
}
export default function Home() {
  //TODO: Redux
  const [videoUrl, setVideoUrl] = useState<string | undefined>();
  const [musicBrainzId, setMusicBrainzId] = useState<string | undefined>();
  const [status, setStatus] = useState<STATUS>(STATUS.READY);
  async function downloadVideo() {
    if (!videoUrl) return;
    try {
      setStatus(STATUS.DOWNLOADING);
      const videoID = YouTubeVideoId(videoUrl);
      const backendURL = `http://localhost:3333/download?id=${videoID}&musicbrainzid=${musicBrainzId}`;
      const response = await fetch(backendURL);
      if (!response.ok) throw new Error(`HTTP error: ${response.status}`);

      setStatus(STATUS.FINISHED);

    } catch (err) {
      console.error("Download failed:", err);
      setStatus(STATUS.ERROR);
    }
  }

  function getStatusText() {
    switch (status) {
      case STATUS.READY: return <Button variant="contained" startIcon={<Download />} onClick={() => downloadVideo()}>
        Runnerladen
      </Button>;
      case STATUS.DOWNLOADING: return <CircularProgress color="secondary" />;
      case STATUS.FINISHED: return (
        <div>
          <Typography>Download successfull! 👺</Typography>
          <Button variant="contained" startIcon={<Download />} onClick={() => downloadVideo()}>Download</Button>
        </div>
      )
      case STATUS.ERROR: return (
        <>
          <Typography>Error downloading 🥴</Typography>
          <Button variant="contained" startIcon={<Download />} onClick={() => downloadVideo()}>
            Try again 🤠
          </Button>
        </>
      )
    }
  }
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <div className="flex min-h-screen flex-col p-4">
        <main className="flex flex-col min-h-screen w-full items-center">
          <Typography color="textPrimary" variant="h3">Audio Grabbler</Typography>
            <VideoForm setURL={(e: string) => setVideoUrl(e)} setMusicBrainzId={(e: string) => setMusicBrainzId(e)}></VideoForm>
            <YoutubeVideo videoURL={videoUrl}></YoutubeVideo>
            <div className="pt-4">
              {getStatusText()}
            </div>
        </main>
      </div>
    </ThemeProvider>
  );
}
