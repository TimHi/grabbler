import { createTheme, ThemeOptions } from "@mui/material";

const themeOptions: ThemeOptions = {
  palette: {
    primary: {
      main: "#38bdf8"
    },
    secondary: {
      main: "#fb7185"
    },
    background: {
      default: "#f5f7fb",
      paper: "#eef2f8"
    },
    text: {
      primary: "#0b1220",
      secondary: "#334155"
    }
  },
  shape: {
    borderRadius: 20
  },
  typography: {
    fontFamily: "\"Space Grotesk\", \"Roboto\", system-ui, -apple-system, sans-serif",
    h3: {
      fontWeight: 700,
      letterSpacing: "-0.02em"
    },
    button: {
      fontWeight: 600,
      textTransform: "none",
      letterSpacing: "0.01em"
    }
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 999,
          paddingLeft: "28px",
          paddingRight: "28px",
          paddingTop: "10px",
          paddingBottom: "10px",
          transition: "box-shadow 200ms ease"
        },
        contained: {
          background: "linear-gradient(90deg, rgba(56,189,248,1) 0%, rgba(236,72,153,1) 100%)",
          boxShadow: "0 18px 45px rgba(56,189,248,0.35)",
          "&:hover": {
            background: "linear-gradient(90deg, rgba(56,189,248,1) 0%, rgba(236,72,153,1) 100%)",
            boxShadow: "0 18px 45px rgba(56,189,248,0.45)"
          }
        }
      }
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          borderRadius: 20,
          backgroundColor: "rgba(255, 255, 255, 0.7)",
          boxShadow: "0 12px 30px rgba(15, 23, 42, 0.12)",
          backdropFilter: "blur(10px)",
          "&:hover .MuiOutlinedInput-notchedOutline": {
            borderColor: "rgba(100, 116, 139, 0.6)"
          },
          "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
            borderColor: "rgba(0, 122, 255, 0.6)",
            boxShadow: "0 0 0 2px rgba(0, 122, 255, 0.2)"
          }
        },
        notchedOutline: {
          borderColor: "rgba(100, 116, 139, 0.35)"
        },
        input: {
          color: "rgba(15, 23, 42, 0.9)",
          paddingTop: "16px",
          paddingBottom: "16px"
        }
      }
    },
    MuiInputLabel: {
      styleOverrides: {
        root: {
          color: "rgba(71, 85, 105, 0.95)",
          "&.Mui-focused": {
            color: "rgba(30, 64, 175, 0.95)"
          }
        }
      }
    }
  }
};

const theme = createTheme(themeOptions);
export default theme;
