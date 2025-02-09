import React, { useState, useContext } from 'react';
import config from '../configs/config';
import { Button } from '@heroui/button';
import { Progress } from '@heroui/progress';
import { uploadChunk, createChunk, chunkAndUpload } from '../services/chunkFile';

import { FileContext } from '../contexts/FileContext';

const FileUploader = () => {
  const { file, setFile, progress, setProgress } = useContext(FileContext);
  const uploadFile = async () => {
    if (!file) return;
    chunkAndUpload(setProgress, file);
}
  return (
    <>
    <Progress isStriped aria-label="Loading..." className="w-[70%]" color="secondary" value={progress}/>
    <Button
      color="primary"
      isDisabled={!file}
      className="text-md w-[80%]"
      size="lg"
      onPress={uploadFile}
    >
      Send Files
    </Button>
    </>
  )
}

export default FileUploader;