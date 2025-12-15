import CloseRoundedIcon from "@mui/icons-material/CloseRounded";
import {
  Alert,
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { ChangeEvent, useCallback, useEffect, useMemo, useState } from "react";
import {
  LoadSystemConfig,
  SaveSystemConfig,
} from "../../wailsjs/go/configs/ConfigService";
import type { configs } from "../../wailsjs/go/models";

type SystemConfig = configs.System;

type AnchorPosition = { left: number; top: number };

type SystemConfigDialogProps = {
  open: boolean;
  onClose: () => void;
  anchorPosition?: AnchorPosition | null;
};

const getErrorMessage = (error: unknown) => {
  if (error instanceof Error) {
    return error.message;
  }
  if (typeof error === "string") {
    return error;
  }
  return "Something went wrong.";
};

const fields: {
  key: keyof SystemConfig;
  label: string;
  helperText?: string;
  multiline?: boolean;
  minRows?: number;
}[] = [
  {
    key: "udpDBDir",
    label: "UDP database path",
    helperText: "Path to udp.db file used for pet data.",
  },
  {
    key: "staticAssetsDBDir",
    label: "Static assets database path",
    helperText: "Path to static_assets.db used for items and activities.",
  },
  {
    key: "staticAssetsSQLDir",
    label: "Static assets SQL path",
    helperText: "SQL file that seeds the static assets database.",
  },
  {
    key: "petImageFolder",
    label: "Pet images folder",
    helperText: "Folder under assets/petImages to use for the pet.",
  },
  {
    key: "itemsImageFolder",
    label: "Items images folder",
    helperText: "Folder under assets/itemImages to use for items.",
  },
  {
    key: "activitiesImageFolder",
    label: "Activities images folder",
    helperText: "Folder under assets/activityImages to use for activities.",
  },
  {
    key: "chatRolePlayContext",
    label: "Chat role-play context",
    helperText: "Long-form prompt sent to the pet chat model.",
    multiline: true,
    minRows: 6,
  },
];

export function SystemConfigDialog({
  open,
  onClose,
  anchorPosition,
}: SystemConfigDialogProps) {
  const [draft, setDraft] = useState<SystemConfig | null>(null);
  const [savedConfig, setSavedConfig] = useState<SystemConfig | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [saveError, setSaveError] = useState<string | null>(null);
  const [saveMessage, setSaveMessage] = useState<string | null>(null);

  useEffect(() => {
    if (!open) {
      return;
    }

    let cancelled = false;

    const loadConfig = async () => {
      setIsLoading(true);
      setLoadError(null);
      setSaveMessage(null);
      setSaveError(null);

      try {
        const config = await LoadSystemConfig();
        if (cancelled) {
          return;
        }
        setDraft(config);
        setSavedConfig(config);
      } catch (error) {
        if (!cancelled) {
          const message = getErrorMessage(error);
          setLoadError(message);
          setDraft(null);
          setSavedConfig(null);
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    };

    void loadConfig();

    return () => {
      cancelled = true;
    };
  }, [open]);

  const handleChange =
    (key: keyof SystemConfig) =>
    (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const value = event.target.value;
      setDraft((previous) =>
        previous
          ? {
              ...previous,
              [key]: value,
            }
          : previous
      );
    };

  const handleSave = useCallback(async () => {
    if (!draft) {
      return;
    }
    setIsSaving(true);
    setSaveError(null);
    setSaveMessage(null);
    try {
      await SaveSystemConfig(draft);
      setSavedConfig(draft);
      setSaveMessage("Saved. Restart the app to apply these settings.");
    } catch (error) {
      const message = getErrorMessage(error);
      setSaveError(message);
    } finally {
      setIsSaving(false);
    }
  }, [draft]);

  const anchoredDialogSx = anchorPosition
    ? {
        "& .MuiDialog-container": {
          alignItems: "flex-start",
          justifyContent: "flex-start",
        },
        "& .MuiPaper-root": {
          margin: 0,
          position: "absolute",
          left: anchorPosition.left,
          top: anchorPosition.top,
          transform: "none",
        },
      }
    : undefined;

  const isDirty = useMemo(() => {
    if (!draft || !savedConfig) {
      return false;
    }
    return fields.some(({ key }) => draft[key] !== savedConfig[key]);
  }, [draft, savedConfig]);

  const isActionDisabled = !draft || isSaving || isLoading || !isDirty;

  return (
    <Dialog
      open={open}
      onClose={onClose}
      fullWidth
      maxWidth="md"
      keepMounted
      sx={anchoredDialogSx}
      slotProps={{
        backdrop: {
          sx: {
            backgroundColor: "transparent",
          },
        },
        paper: {
          sx: {
            "--wails-draggable": "no-drag",
            WebkitAppRegion: "no-drag",
            borderRadius: 2,
            backgroundColor: "rgba(24,24,24,0.95)",
            color: "#f5f5f5",
            boxShadow: 8,
            px: 1,
          },
        },
      }}
    >
      <DialogTitle
        sx={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          pr: 1,
        }}
      >
        System settings
        <IconButton
          aria-label="Close settings"
          onClick={onClose}
          size="small"
          sx={{ color: "#f5f5f5" }}
        >
          <CloseRoundedIcon fontSize="small" />
        </IconButton>
      </DialogTitle>

      <DialogContent dividers sx={{ minHeight: 360 }}>
        <Stack spacing={1.25}>
          {loadError ? <Alert severity="error">{loadError}</Alert> : null}
          {saveError ? <Alert severity="error">{saveError}</Alert> : null}
          {saveMessage ? <Alert severity="success">{saveMessage}</Alert> : null}

          {isLoading ? (
            <Box sx={{ display: "flex", justifyContent: "center", py: 6 }}>
              <CircularProgress size={28} />
            </Box>
          ) : !draft ? (
            <Typography variant="body2" sx={{ color: "#f5f5f5" }}>
              Unable to load config.
            </Typography>
          ) : (
            <Stack spacing={1.25}>
              {fields.map((field) => (
                <TextField
                  key={field.key}
                  label={field.label}
                  value={draft[field.key] ?? ""}
                  onChange={handleChange(field.key)}
                  fullWidth
                  size="small"
                  multiline={field.multiline}
                  minRows={field.minRows}
                  disabled={isSaving}
                  sx={{
                    "& .MuiFormHelperText-root": {
                      color: "rgba(255,255,255,0.6)",
                    },
                    "& .MuiInputBase-root": {
                      backgroundColor: "rgba(0,0,0,0.35)",
                      color: "#f5f5f5",
                      borderRadius: 1.5,
                    },
                    "& .MuiOutlinedInput-notchedOutline": {
                      borderColor: "rgba(255,255,255,0.2)",
                    },
                  }}
                  slotProps={{
                    inputLabel: { sx: { color: "rgba(255,255,255,0.8)" } },
                  }}
                  helperText={field.helperText}
                />
              ))}

              <Box
                sx={{
                  display: "flex",
                  justifyContent: "flex-end",
                  gap: 1,
                  mt: 0.5,
                }}
              >
                <Button
                  variant="text"
                  onClick={onClose}
                  sx={{ color: "rgba(255,255,255,0.78)" }}
                >
                  Close
                </Button>
                <Button
                  variant="contained"
                  disabled={isActionDisabled}
                  onClick={() => void handleSave()}
                  sx={{ px: 3 }}
                >
                  {isSaving ? "Saving..." : "Save"}
                </Button>
              </Box>
            </Stack>
          )}
        </Stack>
      </DialogContent>
    </Dialog>
  );
}
