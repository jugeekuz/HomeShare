import React from 'react'
import { Spinner } from '@heroui/react';
const LoadingPage: React.FC = () => {
    return (
        <div className="flex justify-center items-center">
            <Spinner classNames={{label: "text-gray-50 mt-4 font-brsonoma"}} label="We are trying to authenticate you" variant="spinner" />
        </div>
    )
}

export default LoadingPage;