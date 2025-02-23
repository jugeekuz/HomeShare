import React from 'react';
import { Button } from '@heroui/button';
import { Progress } from '@heroui/progress';

import { useFileContext } from '../contexts/FileContext.tsx';

const FileUploader : React.FC = () => {
    const { files, progress, uploadFiles } = useFileContext();

    return (
        <>
        { progress !== 0 && (
            <Progress 
            isStriped 
            aria-label="Loading..." 
            className="w-[70%]" 
            color="secondary" 
            value={progress}
            />
        )}
        <Button
            color="primary"
            isDisabled={!files || Object.keys(files).length === 0}
            className="text-md w-[80%]"
            size="lg"
            onPress={uploadFiles}
        >
            Send Files
        </Button>
        </>
    )
}

export default FileUploader;