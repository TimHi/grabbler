'use client'

import YoutubeVideo from "./components/YoutubeVideo";
import { useState } from "react";
import VideoForm from "./components/VideoForm";
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme';
import { Typography } from "@mui/material";


export default function Home() {
  const [videoUrl, setVideoUrl] = useState<string | undefined>();
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <div className="flex min-h-screen items-center flex-col p-4">
        <main className="min-h-screen w-full p-4 items-center">
          <VideoForm setURL={(e: string) => setVideoUrl(e)}></VideoForm>
          <YoutubeVideo videoURL={videoUrl}></YoutubeVideo> 
          <Typography>Video URL rein und ab gehts</Typography>
        </main>
      </div>
    </ThemeProvider>
  );
}
