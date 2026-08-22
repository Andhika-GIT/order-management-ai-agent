"use client";

import { IconBell, IconBellFilled, IconX } from "@tabler/icons-react";
import { useEffect, useRef, useState } from "react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";

interface Notification {
  id: string;
  message: string;
  timestamp: Date;
}

export const Notification = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | undefined>(undefined);

  useEffect(() => {
    // Load from localStorage
    const saved = localStorage.getItem("notifications");
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        setNotifications(parsed.map((n: any) => ({
          ...n,
          timestamp: new Date(n.timestamp)
        })));
      } catch (e) {
        console.error("Failed to load notifications", e);
      }
    }
  }, []);

  useEffect(() => {
    // Save to localStorage
    if (notifications.length > 0) {
      localStorage.setItem("notifications", JSON.stringify(notifications));
    }
  }, [notifications]);

  useEffect(() => {
    const connect = () => {
      const ws = new WebSocket("ws://localhost:9080");
      wsRef.current = ws;

      ws.onopen = () => {
        console.log("WebSocket Connected!");
        setIsConnected(true);
      };

      ws.onmessage = (event) => {
        console.log("Message received:", event.data);
        
        const newNotification: Notification = {
          id: Date.now().toString(),
          message: event.data,
          timestamp: new Date(),
        };

        setNotifications((prev) => {
          const updated = [newNotification, ...prev];
          // Keep only last 50
          return updated.slice(0, 50);
        });
      };

      ws.onerror = (error) => {
        console.error("WebSocket Error:", error);
      };

      ws.onclose = () => {
        console.log("WebSocket Disconnected");
        setIsConnected(false);
        // Retry after 3 seconds
        reconnectTimeoutRef.current = setTimeout(connect, 3000);
      };
    };

    connect();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  const clearAll = () => {
    setNotifications([]);
    localStorage.removeItem("notifications");
  };

  const removeNotification = (id: string) => {
    setNotifications((prev) => prev.filter((n) => n.id !== id));
  };

  const formatTime = (date: Date) => {
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    
    if (diffMins < 1) return "Just now";
    if (diffMins < 60) return `${diffMins}m ago`;
    
    const diffHours = Math.floor(diffMins / 60);
    if (diffHours < 24) return `${diffHours}h ago`;
    
    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays}d ago`;
  };

  return (
    <Popover>
      <PopoverTrigger asChild>
        <button className="relative focus:outline-none">
          {notifications.length > 0 ? (
            <>
              <IconBellFilled className="text-red-500" size={28} />
              <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full min-w-5 h-5 px-1 flex items-center justify-center font-medium">
                {notifications.length > 99 ? "99+" : notifications.length}
              </span>
            </>
          ) : (
            <IconBell size={28} />
          )}
        </button>
      </PopoverTrigger>
      <PopoverContent className="w-96 p-0" align="end">
        <div className="flex items-center justify-between p-4 border-b">
          <div>
            <h3 className="font-semibold">Notifications</h3>
            <div className="flex items-center gap-2 mt-1">
              <div
                className={`w-2 h-2 rounded-full ${
                  isConnected ? "bg-green-500" : "bg-red-500"
                }`}
              />
              <span className="text-xs text-muted-foreground">
                {isConnected ? "Connected" : "Reconnecting..."}
              </span>
            </div>
          </div>
          {notifications.length > 0 && (
            <Button
              variant="ghost"
              size="sm"
              onClick={clearAll}
              className="text-xs"
            >
              Clear all
            </Button>
          )}
        </div>

        <div className="max-h-[400px] overflow-y-auto">
          {notifications.length === 0 ? (
            <div className="p-8 text-center text-muted-foreground">
              <IconBell size={48} className="mx-auto mb-2 opacity-20" />
              <p className="text-sm">No notifications yet</p>
            </div>
          ) : (
            <div className="divide-y">
              {notifications.map((notif) => (
                <div
                  key={notif.id}
                  className="p-3 hover:bg-accent/50 transition-colors group"
                >
                  <div className="flex items-start justify-between gap-2">
                    <div className="flex-1 min-w-0">
                      <p className="text-sm break-words">{notif.message}</p>
                      <p className="text-xs text-muted-foreground mt-1">
                        {formatTime(notif.timestamp)}
                      </p>
                    </div>
                    <button
                      onClick={() => removeNotification(notif.id)}
                      className="opacity-0 group-hover:opacity-100 transition-opacity p-1 hover:bg-accent rounded"
                    >
                      <IconX size={16} />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
};