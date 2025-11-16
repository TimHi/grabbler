import { Typography } from "@mui/material";
import TextField from "@mui/material/TextField";

export interface VideoProps {
    setURL: (value: string) => void
}

export default function VideoForm(props: VideoProps) {
    return (
        <div className="flex flex-col p-4 items-center min-w-full gap-4">
            <Typography color="textPrimary" variant="h3">Audio Grabbler</Typography>
            <TextField sx={{
                minWidth: 500
            }} id="outlined-basic" onChange={(v) => props.setURL(v.target.value)} variant="outlined" placeholder="youtube.com/bla" label="Video URL" />
        </div>
    );
}