"use client";

import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Dropzone,
  DropzoneContent,
  DropzoneEmptyState,
} from "@/components/ui/shadcn-io/dropzone";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { confirmUploadOrderAttachment, getOrderImagePresignedUrl, uploadFileToPresignedUrl } from "@/app/action/order";
import { Error } from "@/lib/types";

type UploadAttachmentModalProps = {
  orderId: number;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export const UploadAttachmentModal: React.FC<UploadAttachmentModalProps> = ({
  orderId,
  open,
  onOpenChange,
}) => {
  const [files, setFiles] = useState<File[] | undefined>();
  const [isLoading, setIsLoading] = useState(false);

  const handleDrop = (accepted: File[]) => {
    setFiles(accepted);
  };

  const closeModal = () => {
    setFiles(undefined);
    onOpenChange(false);
  };

  const handleUpload = async () => {
    const file = files?.[0];
    if (!file) return;

    setIsLoading(true);
    try {
      const presigned = await getOrderImagePresignedUrl(
        orderId,
        file.name,
        file.type
      );

      if (!presigned) {
        throw { success: false, code: 500, message: "Failed to get presigned URL" } as Error;
      }

      await uploadFileToPresignedUrl(presigned.url, file);
      await confirmUploadOrderAttachment(orderId, presigned.key)

      // TODO: persist `presigned.key` against the order
      toast.success("File uploaded successfully");
      closeModal();
    } catch (e) {
      const error = e as Error;
      toast.error(error.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        if (!next) setFiles(undefined);
        onOpenChange(next);
      }}
    >
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Upload Attachment</DialogTitle>
        </DialogHeader>
        <Dropzone
          accept={{
            "image/png": [".png"],
            "image/jpeg": [".jpg", ".jpeg"],
          }}
          maxSize={1024 * 1024 * 10}
          onDrop={handleDrop}
          onError={(e) => toast.error(e.message)}
          src={files}
        >
          <DropzoneEmptyState />
          <DropzoneContent />
        </Dropzone>
        <DialogFooter>
          <Button disabled={!files || isLoading} onClick={handleUpload}>
            {isLoading ? "Please wait..." : "Upload"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
