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
  Typography,
} from "@mui/material";
import { useCallback, useEffect, useMemo, useState } from "react";
import {
  LoadAllItems,
  LoadItemFrameByID,
} from "../../../wailsjs/go/items/ItemsMeta";
import { GetPetStatus } from "../../../wailsjs/go/pet/PetMeta";

type PetStatus = Awaited<ReturnType<typeof GetPetStatus>>;

type AnchorPosition = { left: number; top: number };

export type ShopItemAttributes = {
  experience: number;
  water: number;
  hunger: number;
  health: number;
  mood: number;
  energy: number;
  money: number;
};

export type ShopItem = {
  id: number;
  path: string;
  name: string;
  type: string;
  attributes: ShopItemAttributes;
  description: string;
};

type PetShopDialogProps = {
  open: boolean;
  onClose: () => void;
  onPurchaseItem: (item: ShopItem, itemFrame: string | null) => Promise<void>;
  anchorPosition?: AnchorPosition | null;
};

const formatMoney = (value: number | null) => {
  if (typeof value !== "number") {
    return "—";
  }
  return value.toLocaleString();
};

const formatAmount = (value: number) => value.toLocaleString();

const getItemMoneyDelta = (item: ShopItem) => {
  const nestedMoney = item.attributes?.money;
  if (typeof nestedMoney === "number") {
    return nestedMoney;
  }
  const directMoney = (item as unknown as Record<string, unknown>)["money"];
  if (typeof directMoney === "number") {
    return directMoney;
  }
  return 0;
};

const getItemPriceLabel = (moneyDelta: number) => {
  if (moneyDelta === 0) {
    return "Free";
  }
  if (moneyDelta < 0) {
    return `Cost $${formatAmount(-moneyDelta)}`;
  }
  return `Gives $${formatAmount(moneyDelta)}`;
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

export function PetShopDialog({
  open,
  onClose,
  onPurchaseItem,
  anchorPosition,
}: PetShopDialogProps) {
  const [items, setItems] = useState<ShopItem[]>([]);
  const [petStatus, setPetStatus] = useState<PetStatus | null>(null);
  const [itemFrames, setItemFrames] = useState<Record<number, string>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [isBuyingItemId, setIsBuyingItemId] = useState<number | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [actionMessage, setActionMessage] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

  const refreshPetStatus = useCallback(async () => {
    try {
      const snapshot = await GetPetStatus();
      setPetStatus(snapshot);
    } catch (error) {
      console.error("GetPetStatus failed", error);
    }
  }, []);

  useEffect(() => {
    if (!open) {
      return;
    }

    let cancelled = false;

    const loadShop = async () => {
      setIsLoading(true);
      setLoadError(null);
      setActionMessage(null);
      setActionError(null);

      try {
        const [itemsSnapshot] = await Promise.all([LoadAllItems()]);
        if (cancelled) {
          return;
        }
        setItems((itemsSnapshot ?? []) as unknown as ShopItem[]);
        void refreshPetStatus();
      } catch (error) {
        if (cancelled) {
          return;
        }
        console.error("LoadAllItems failed", error);
        setLoadError("Shop inventory unavailable.");
        setItems([]);
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    };

    void loadShop();

    return () => {
      cancelled = true;
    };
  }, [open, refreshPetStatus]);

  useEffect(() => {
    if (!open || items.length === 0) {
      return;
    }

    let cancelled = false;

    const loadFrames = async () => {
      const entries = await Promise.all(
        items.map(async (item) => {
          try {
            const frame = await LoadItemFrameByID(item.id);
            return [item.id, frame] as const;
          } catch (error) {
            console.warn("LoadItemFrameByID failed", item.id, error);
            return [item.id, null] as const;
          }
        })
      );

      if (cancelled) {
        return;
      }

      setItemFrames((previous) => {
        const next: Record<number, string> = { ...previous };
        entries.forEach(([id, frame]) => {
          if (frame) {
            next[id] = frame;
          }
        });
        return next;
      });
    };

    void loadFrames();

    return () => {
      cancelled = true;
    };
  }, [items, open]);

  const handlePurchase = useCallback(
    async (item: ShopItem) => {
      if (isBuyingItemId !== null) {
        return;
      }

      setIsBuyingItemId(item.id);
      setActionMessage(null);
      setActionError(null);

      let frame: string | null = itemFrames[item.id] ?? null;
      if (!frame) {
        try {
          frame = await LoadItemFrameByID(item.id);
          setItemFrames((previous) => ({
            ...previous,
            [item.id]: frame ?? "",
          }));
        } catch (error) {
          console.warn(
            "LoadItemFrameByID failed during purchase",
            item.id,
            error
          );
          frame = null;
        }
      }

      try {
        await onPurchaseItem(item, frame);
        setActionMessage(`Used ${item.name}.`);
        await refreshPetStatus();
      } catch (error) {
        const message = getErrorMessage(error);
        setActionError(message);
      } finally {
        setIsBuyingItemId(null);
      }
    },
    [isBuyingItemId, itemFrames, onPurchaseItem, refreshPetStatus]
  );

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

  const money = typeof petStatus?.money === "number" ? petStatus.money : null;

  const canPurchase = useMemo(() => {
    if (money === null) {
      return () => false;
    }
    return (item: ShopItem) => money + getItemMoneyDelta(item) >= 0;
  }, [money]);

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
        <Box sx={{ display: "flex", flexDirection: "column" }}>
          <Typography component="span" sx={{ fontWeight: 700 }}>
            Shop
          </Typography>
          <Typography
            component="span"
            variant="caption"
            sx={{ color: "rgba(255,255,255,0.65)" }}
          >
            Money ${formatMoney(money)}
          </Typography>
        </Box>
        <IconButton
          aria-label="Close shop"
          onClick={onClose}
          size="small"
          sx={{ color: "#f5f5f5" }}
        >
          <CloseRoundedIcon fontSize="small" />
        </IconButton>
      </DialogTitle>

      <DialogContent dividers sx={{ minHeight: 340 }}>
        <Stack spacing={1.5}>
          {loadError ? <Alert severity="error">{loadError}</Alert> : null}
          {actionError ? <Alert severity="error">{actionError}</Alert> : null}
          {actionMessage ? (
            <Alert severity="success">{actionMessage}</Alert>
          ) : null}

          {isLoading ? (
            <Box sx={{ display: "flex", justifyContent: "center", py: 6 }}>
              <CircularProgress size={28} />
            </Box>
          ) : items.length === 0 ? (
            <Typography variant="body2" sx={{ color: "rgba(255,255,255,0.7)" }}>
              No items available.
            </Typography>
          ) : (
            <Stack spacing={1.25}>
              {items.map((item) => {
                const moneyDelta = getItemMoneyDelta(item);
                const isBusy = isBuyingItemId === item.id;
                const isDisabled = !canPurchase(item) || isBusy;
                const frame = itemFrames[item.id];

                return (
                  <Box
                    key={item.id}
                    sx={{
                      display: "flex",
                      gap: 1.5,
                      alignItems: "flex-start",
                      borderRadius: 2,
                      border: "1px solid rgba(255,255,255,0.12)",
                      backgroundColor: "rgba(0,0,0,0.28)",
                      px: 1.5,
                      py: 1.25,
                    }}
                  >
                    <Box
                      sx={{
                        width: 56,
                        height: 56,
                        borderRadius: 2,
                        backgroundColor: "rgba(255,255,255,0.06)",
                        border: "1px solid rgba(255,255,255,0.08)",
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        overflow: "hidden",
                        flexShrink: 0,
                      }}
                    >
                      {frame ? (
                        <Box
                          component="img"
                          src={frame}
                          alt={item.name}
                          sx={{
                            width: "100%",
                            height: "100%",
                            objectFit: "contain",
                            imageRendering: "pixelated",
                          }}
                        />
                      ) : (
                        <Typography
                          variant="caption"
                          sx={{ color: "rgba(255,255,255,0.55)" }}
                        >
                          ...
                        </Typography>
                      )}
                    </Box>

                    <Box sx={{ flex: 1, minWidth: 0 }}>
                      <Box
                        sx={{
                          display: "flex",
                          justifyContent: "space-between",
                          gap: 1,
                        }}
                      >
                        <Box sx={{ minWidth: 0 }}>
                          <Typography
                            variant="subtitle2"
                            sx={{ fontWeight: 700, lineHeight: 1.2 }}
                            noWrap
                          >
                            {item.name}
                          </Typography>
                          <Typography
                            variant="caption"
                            sx={{ color: "rgba(255,255,255,0.6)" }}
                          >
                            {getItemPriceLabel(moneyDelta)}
                          </Typography>
                        </Box>

                        <Button
                          variant="contained"
                          size="small"
                          disabled={isDisabled}
                          onClick={() => void handlePurchase(item)}
                          sx={{ alignSelf: "flex-start", px: 2.25 }}
                        >
                          {isBusy ? "..." : "Buy & use"}
                        </Button>
                      </Box>

                      {item.description ? (
                        <Typography
                          variant="body2"
                          sx={{
                            mt: 0.75,
                            color: "rgba(255,255,255,0.75)",
                            lineHeight: 1.35,
                          }}
                        >
                          {item.description}
                        </Typography>
                      ) : null}

                      {!canPurchase(item) && money !== null ? (
                        <Typography
                          variant="caption"
                          sx={{ mt: 0.75, display: "block", color: "#fca5a5" }}
                        >
                          Not enough money.
                        </Typography>
                      ) : null}
                    </Box>
                  </Box>
                );
              })}
            </Stack>
          )}
        </Stack>
      </DialogContent>
    </Dialog>
  );
}
