import { Download } from "@mui/icons-material";
import { Button, CircularProgress, Typography } from "@mui/material";
import TextField from "@mui/material/TextField";
import { useState } from "react";


export interface VideoProps {
    setURL: (value: string) => void,
    setMusicBrainzId: (value: string) => void
}

export default function VideoForm(props: VideoProps) {
    
    return (
        <div className="flex flex-col p-4 items-center min-w-full gap-4">
            <Typography color="textPrimary" variant="h3">Audio Grabbler</Typography>
            <TextField sx={{
                minWidth: 500
            }} id="outlined-basic" onChange={(v) => props.setURL(v.target.value)} variant="outlined" placeholder="youtube.com/bla" label="Video URL" />
            <TextField sx={{
                                minWidth: 500
                            }} id="outlined-basic" onChange={(v) => props.setMusicBrainzId(v.target.value)} variant="outlined" placeholder="e.g. 12a5b094-3804-4c97-82b8-9c7cc5d4f4ab" label="Musicbrainz Track ID" />
        </div>
    );
}