import React from 'react';
import { Button } from '@heroui/button';
import { Spinner } from '@heroui/spinner';
import { Progress } from '@heroui/progress';

import { useFileContext } from '../contexts/FileContext.tsx';

const FileUploader : React.FC = () => {
    const { files, filesReady, progress, filesUploading, uploadFiles } = useFileContext();

    return (
        <>
        {/* { progress !== 0 && (
            <Progress 
            isStriped 
            aria-label="Loading..." 
            className="w-[70%]" 
            color="secondary" 
            value={progress}
            />
        )} */}
        <div className="w-full px-2">
            <Button
                color="primary"
                isDisabled={!files || Object.keys(files).length === 0 || !filesReady}
                className="text-md w-full bg-primary-gradient"
                size="md"
                radius="sm"
                onPress={uploadFiles}
            >
                {((files && !filesReady) || filesUploading) ? <Spinner color="default"/> : "Send Files"}
            </Button>
        </div>
        </>
    )
}

export default FileUploader;