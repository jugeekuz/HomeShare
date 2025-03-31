import React, { createContext, useContext, ReactNode } from "react";

import { addToast, ToastProvider } from "@heroui/react";

interface NotificationContextType {
    notifyInfo: (title: string, message: string) => void;
    notifySuccess: (title: string, message: string) => void;
    notifyError: (title: string, message: string) => void;
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined);

export const useNotificationContext = () => {
    const context = useContext(NotificationContext);
    if (!context) {
        throw new Error("useNotificationContext must be used within a NotificationProvider");
    }
    return context;
};

export const NotificationProvider: React.FC<{ children: ReactNode }> = ({ children }) => {

    const notifyInfo = (title: string, message: string) => addToast({
        title: title,
        color: "default",
        description: message,
    });
    const notifySuccess = (title: string, message: string) => addToast({
        title: title,
        color: "success",
        description: message,
    });
    const notifyError = (title: string, message: string) => addToast({
        title: title,
        color: "danger",
        description: message,
    });
  
    return (
        <NotificationContext.Provider value={{ notifyInfo, notifySuccess, notifyError }}>
            <ToastProvider placement="top-center" toastOffset={60} />
                {children}
        </NotificationContext.Provider>
    );
};
