import { createTheme, ThemeOptions } from "@mui/material";


const themeOptions: ThemeOptions = {
  palette: {
    primary: {
      main: '#b73856',
    },
    secondary: {
      main: '#f50057',
    },
    background: {
      default: '#272139',
      paper: '#2b2a2a',
    },
    text: {
        primary: '#b73856',
        secondary: '#b73856',
    }
  },
};

const theme = createTheme(themeOptions);
export default theme;