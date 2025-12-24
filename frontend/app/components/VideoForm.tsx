import TextField from "@mui/material/TextField";


export interface VideoProps {
    setURL: (value: string) => void,
    setMusicBrainzId: (value: string) => void
}

export default function VideoForm(props: VideoProps) {
    
    return (
        <div className="video-form">
            <TextField id="outlined-basic" onChange={(v) => props.setURL(v.target.value)} variant="outlined" placeholder="youtube.com/bla" label="Video URL" />
            <TextField id="outlined-basic" onChange={(v) => props.setMusicBrainzId(v.target.value)} variant="outlined" placeholder="e.g. 12a5b094-3804-4c97-82b8-9c7cc5d4f4ab" label="Musicbrainz Track ID" />
        </div>
    );
}
