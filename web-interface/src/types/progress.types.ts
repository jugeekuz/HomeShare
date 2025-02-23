
export interface ProgressBarProps {
    className?: string;
}
  
export interface ProgressBarRef {
    updateProgress: (newProgress: number) => void;
}

export type UpdateProgress = (newProgress: number) => void

export type ProgressBarRefs = Record<string, ProgressBarRef | null>;