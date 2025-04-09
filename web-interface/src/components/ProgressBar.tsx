import { useState, forwardRef, useImperativeHandle }  from 'react'
import { Progress } from '@heroui/progress';

import { ProgressBarProps, ProgressBarRef } from '../types';

const ProgressBar = forwardRef<ProgressBarRef, ProgressBarProps>((props, ref) => {
    const [progress, setProgress] = useState(0);
  
    useImperativeHandle(ref, () => ({
      updateProgress: (newProgress: number) => {
        setProgress(newProgress);
      },
      getProgress: () => progress,
    }), []);
  
    return (
        <>
        <Progress  
            size={props.size ? props.size : "md"}
            aria-label="Loading..." 
            classNames={{
                track: "bg-white border border-gray-200",
                indicator: "bg-blue-500"
              }}
            className={props.className}
            color="primary" 
            value={progress}
        />
        <span className="text-xs text-gray-500 ml-2">
            {`${progress}%`}
        </span>
        </>
    );
});

export default ProgressBar
