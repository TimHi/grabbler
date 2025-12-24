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
      default: "#0f172a",
      paper: "#0b1120"
    },
    text: {
      primary: "#f8fafc",
      secondary: "#cbd5f5"
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
          paddingBottom: "10px"
        },
        contained: {
          background: "linear-gradient(90deg, rgba(56,189,248,1) 0%, rgba(236,72,153,1) 100%)",
          boxShadow: "0 18px 45px rgba(56,189,248,0.35)",
          "&:hover": {
            background: "linear-gradient(90deg, rgba(14,165,233,1) 0%, rgba(248,113,113,1) 100%)"
          }
        }
      }
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          borderRadius: 20,
          backgroundColor: "rgba(15, 23, 42, 0.75)",
          boxShadow: "0 12px 30px rgba(15, 23, 42, 0.35)",
          backdropFilter: "blur(10px)",
          "&:hover .MuiOutlinedInput-notchedOutline": {
            borderColor: "rgba(148, 163, 184, 0.55)"
          },
          "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
            borderColor: "rgba(56, 189, 248, 0.75)",
            boxShadow: "0 0 0 2px rgba(56, 189, 248, 0.2)"
          }
        },
        notchedOutline: {
          borderColor: "rgba(148, 163, 184, 0.35)"
        },
        input: {
          color: "rgba(248, 250, 252, 0.95)",
          paddingTop: "16px",
          paddingBottom: "16px"
        }
      }
    },
    MuiInputLabel: {
      styleOverrides: {
        root: {
          color: "rgba(148, 163, 184, 0.9)",
          "&.Mui-focused": {
            color: "rgba(125, 211, 252, 0.95)"
          }
        }
      }
    }
  }
};

const theme = createTheme(themeOptions);
export default theme;
