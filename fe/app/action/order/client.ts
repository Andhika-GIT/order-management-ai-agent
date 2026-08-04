"use client"

import { handleFetchResponse, SERVER_BASE_URL_FOR_CLIENT } from "@/lib/helper"
import { Order } from "@/lib/schemas";
import { Error, Paginate } from "@/lib/types";

type FetchResult = Paginate<Order[]>

type PresignResult = {
  url: string;
  key: string;
}

export const getOrderImagePresignedUrl = async (orderId: number, filename: string, contentType: string): Promise<PresignResult | undefined> => {
  const BASE_URL = `${SERVER_BASE_URL_FOR_CLIENT}/upload/order_image/presign`;

  try {
    const response = await fetch(BASE_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        order_id: orderId,
        filename,
        content_type: contentType,
      }),
    });

    return await handleFetchResponse(response);
  } catch (e) {
    throw e as Error;
  }
};

export const uploadFileToPresignedUrl = async (url: string, file: File): Promise<void> => {
  try {
    const response = await fetch(url, {
      method: "PUT",
      headers: { "Content-Type": file.type },
      body: file,
    });

    if (!response.ok) {
      const error: Error = {
        success: false,
        code: response.status,
        message: "Failed to upload file to storage",
      };
      throw error;
    }
  } catch (e) {
    throw e as Error;
  }
};

export const confirmUploadOrderAttachment = async (orderId: number, key: string): Promise<void> => {
  const BASE_URL = `${SERVER_BASE_URL_FOR_CLIENT}/upload/order_image/confirm`;

  try {
    const response = await fetch(BASE_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ order_id: orderId, key }),
    });

    return await handleFetchResponse(response);
  } catch (e) {
    throw e as Error;
  }
};

export const UploadOrderExcel = async (file: File): Promise<string | undefined> => {
  const BASE_URL = `${SERVER_BASE_URL_FOR_CLIENT}/order/upload`;

  try {

    const formData = new FormData()
    formData.append("file", file)

    const response = await fetch(BASE_URL, {
      method: "POST",
      body:formData,
    });

    return await handleFetchResponse(response);
  } catch (e) {
    throw e as Error;
  }
};

export const getAllOrders = async (page: number, pageSize: number, search: null | string = null): Promise<FetchResult | undefined> => {
    const params = new URLSearchParams({
    page: page.toString(),
    per_page: pageSize.toString(),
  });

  if (search && search.trim() !== '') {
    params.append('search', search.trim());
  }

  const BASE_URL = `${SERVER_BASE_URL_FOR_CLIENT}/order?${params.toString()}`;

  try {
    const response = await fetch(BASE_URL, {
      method: "GET",
    });

    return await handleFetchResponse(response);
  } catch (e) {
    throw e as Error;
  }
};