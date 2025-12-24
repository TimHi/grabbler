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
      case STATUS.READY: return (
        <Button variant="contained" startIcon={<Download />} onClick={() => downloadVideo()}>
          Start download
        </Button>
      );
      case STATUS.DOWNLOADING: return <CircularProgress className="status-spinner" />;
      case STATUS.FINISHED: return (
        <div>
          <Typography className="status-text status-text--success">Download complete.</Typography>
          <Button variant="contained" startIcon={<Download />} onClick={() => downloadVideo()}>
            Download again
          </Button>
        </div>
      )
      case STATUS.ERROR: return (
        <>
          <Typography className="status-text status-text--error">Download failed. Please try again.</Typography>
          <Button variant="contained" startIcon={<Download />} onClick={() => downloadVideo()}>
            Retry download
          </Button>
        </>
      )
    }
  }
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <div className="app-shell">
        <div className="app-glow" />
        <main className="app-main">
          <Typography variant="h3" className="app-title">
            Audio Grabber
          </Typography>
          <div className="app-card">
            <div className="flex flex-col gap-6">
              <VideoForm setURL={(e: string) => setVideoUrl(e)} setMusicBrainzId={(e: string) => setMusicBrainzId(e)}></VideoForm>
              <YoutubeVideo videoURL={videoUrl}></YoutubeVideo>
              <div className="flex items-center justify-center pt-2">
                {getStatusText()}
              </div>
            </div>
          </div>
        </main>
      </div>
    </ThemeProvider>
  );
}
