
export interface ProgressBarProps {
    className?: string;
    size?: "sm" | "md" | "lg" | undefined
}
  
export interface ProgressBarRef {
    updateProgress: (newProgress: number) => void;
    getProgress: () => number;
}

export type UpdateProgress = (newProgress: number) => void

export type ProgressBarRefs = Record<string, ProgressBarRef | null>;