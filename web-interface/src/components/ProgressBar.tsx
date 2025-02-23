import React, { useState, forwardRef, useImperativeHandle }  from 'react'
import { Progress } from '@heroui/progress';

import { ProgressBarProps, ProgressBarRef } from '../types';

const ProgressBar = forwardRef<ProgressBarRef, ProgressBarProps>((props, ref) => {
    const [progress, setProgress] = useState(0);
  
    useImperativeHandle(ref, () => ({
      updateProgress: (newProgress: number) => {
        setProgress(newProgress);
      },
    }), []);
  
    return (
        <>
        <Progress 
            isStriped 
            aria-label="Loading..." 
            className={props.className} 
            color="secondary" 
            value={progress}
        />
        </>
    );
});

export default ProgressBar
