import CloseRoundedIcon from "@mui/icons-material/CloseRounded";
import {
  Box,
  Button,
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
  TextField,
  Typography,
} from "@mui/material";
import {
  ChangeEvent,
  FormEvent,
  useEffect,
  useMemo,
  useState,
} from "react";

export type ConversationMessage = {
  id: string;
  role: "user" | "pet";
  text: string;
};

type AnchorPosition = { left: number; top: number };

type PetDialogProps = {
  open: boolean;
  onClose: () => void;
  onSend: (message: string) => void | Promise<void>;
  isBusy?: boolean;
  messages: ConversationMessage[];
  anchorPosition?: AnchorPosition | null;
};

export function PetDialog({
  open,
  onClose,
  onSend,
  isBusy = false,
  messages,
  anchorPosition,
}: PetDialogProps) {
  const [draft, setDraft] = useState("");

  useEffect(() => {
    if (!open) {
      setDraft("");
    }
  }, [open]);

  const isSendDisabled = useMemo(
    () => !draft.trim() || isBusy,
    [draft, isBusy]
  );

  const handleSend = () => {
    const trimmed = draft.trim();
    if (!trimmed || isBusy) {
      return;
    }
    void onSend(trimmed);
    setDraft("");
  };

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

  return (
    <Dialog
      open={open}
      onClose={onClose}
      fullWidth
      maxWidth="sm"
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
        Talk with your pet
        <IconButton
          aria-label="Close dialog"
          onClick={onClose}
          size="small"
          sx={{ color: "#f5f5f5" }}
        >
          <CloseRoundedIcon fontSize="small" />
        </IconButton>
      </DialogTitle>
      <DialogContent
        dividers
        sx={{
          display: "flex",
          flexDirection: "column",
          gap: 2,
          minHeight: 280,
        }}
      >
        <Box
          sx={{
            flex: 1,
            overflowY: "auto",
            display: "flex",
            flexDirection: "column",
            gap: 1.25,
            pr: 1,
          }}
        >
          {messages.length === 0 ? (
            <Typography variant="body2" color="rgba(255,255,255,0.7)">
              Say hi to start a conversation.
            </Typography>
          ) : (
            messages.map((message) => {
              const isUser = message.role === "user";
              return (
                <Box
                  key={message.id}
                  sx={{
                    alignSelf: isUser ? "flex-end" : "flex-start",
                    maxWidth: "80%",
                    backgroundColor: isUser
                      ? "rgba(96, 165, 250, 0.3)"
                      : "rgba(55, 65, 81, 0.85)",
                    borderRadius: 2,
                    px: 1.5,
                    py: 1,
                  }}
                >
                  <Typography
                    variant="caption"
                    sx={{ textTransform: "uppercase", opacity: 0.8 }}
                  >
                    {isUser ? "You" : "Pet"}
                  </Typography>
                  <Typography
                    variant="body2"
                    sx={{ mt: 0.25, lineHeight: 1.4 }}
                  >
                    {message.text}
                  </Typography>
                </Box>
              );
            })
          )}
        </Box>
        <Box
          component="form"
          onSubmit={(event: FormEvent<HTMLFormElement>) => {
            event.preventDefault();
            handleSend();
          }}
          sx={{
            display: "flex",
            gap: 1,
            width: "100%",
          }}
        >
          <TextField
            fullWidth
            size="small"
            placeholder="Tell your pet something..."
            value={draft}
            onChange={(
              event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
            ) => setDraft(event.target.value)}
            disabled={isBusy}
            multiline
            minRows={1}
            maxRows={3}
            sx={{
              "& .MuiInputBase-root": {
                backgroundColor: "rgba(0,0,0,0.35)",
                color: "#f5f5f5",
                borderRadius: 1.5,
              },
              "& .MuiOutlinedInput-notchedOutline": {
                borderColor: "rgba(255,255,255,0.2)",
              },
            }}
          />
          <Button
            type="submit"
            variant="contained"
            disabled={isSendDisabled}
            sx={{ alignSelf: "flex-end", px: 3 }}
          >
            {isBusy ? "..." : "Send"}
          </Button>
        </Box>
      </DialogContent>
    </Dialog>
  );
}
