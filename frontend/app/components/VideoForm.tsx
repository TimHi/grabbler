import TextField from '@mui/material/TextField';
import Tooltip from '@mui/material/Tooltip';

export interface VideoProps {
	setURL: (value: string) => void;
	setMusicBrainzId: (value: string) => void;
}

export default function VideoForm(props: VideoProps) {
	return (
		<div className='video-form'>
			<TextField
				id='outlined-basic'
				onChange={(v) => props.setURL(v.target.value)}
				variant='outlined'
				placeholder='https://www.youtube.com/watch?v=...'
				label='YouTube video URL'
			/>
			<Tooltip title='Use the MusicBrainz recording ID to tag the download with accurate metadata (recording ID, not release ID).'>
				<TextField
					id='outlined-basic'
					onChange={(v) => props.setMusicBrainzId(v.target.value)}
					variant='outlined'
					placeholder='e.g. 12a5b094-3804-4c97-82b8-9c7cc5d4f4ab'
					label='MusicBrainz recording ID'
				/>
			</Tooltip>
		</div>
	);
}
