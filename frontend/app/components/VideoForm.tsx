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
				placeholder='youtube.com/bla'
				label='Video URL'
			/>
			<Tooltip title='The MusicBrainz Recording ID helps to tag the downloaded audio file with correct metadata. Important: recording ID, not release ID!'>
				<TextField
					id='outlined-basic'
					onChange={(v) => props.setMusicBrainzId(v.target.value)}
					variant='outlined'
					placeholder='e.g. 12a5b094-3804-4c97-82b8-9c7cc5d4f4ab'
					label='Musicbrainz Recording ID'
				/>
			</Tooltip>
		</div>
	);
}
