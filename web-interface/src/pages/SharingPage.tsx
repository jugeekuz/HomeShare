import React, {useState} from "react";
import { Tabs, Tab, Card, CardBody, Divider, DropdownMenu } from "@heroui/react";
import { useNavigate } from "react-router-dom";
import UploadComponent from "../components/UploadComponent";
import ShareComponent from "../components/ShareComponent";
import { FaLink } from "react-icons/fa6";
import { FileProvider } from "../contexts/FileContext";
import { LuUpload, LuDownload } from "react-icons/lu";
import DownloadComponent from "../components/DownloadComponent";

type Key = string | number;

const SharingPage: React.FC = () => {
    const [selectedTab, setSelectedTab] = useState<Key>("download")
    // 
    return (
        <div className="flex w-full h-full justify-center items-center">
            {/* Rectangle */}
            <FileProvider>
            <Card className="max-w-full w-[440px]  bg-wprimary">
                <CardBody className="justify-center items-center">
                    <Tabs
                        fullWidth
                        selectedKey={selectedTab}
                        size="sm"
                        aria-label="Tabs colors" 
                        radius="sm" 
                        variant="solid"
                        onSelectionChange={(key) => setSelectedTab(key)}
                        className=""
                        classNames={{
                            base: "bg-wprimary", 
                            tabList: "border border-gray-200 bg-wsecondary p-1 w-[70%] h-9 mx-auto mt-1",
                            tab: "px-4 py-2 text-gray-500 font-normal cursor-pointer transition-all hover:bg-gray-100 data-[selected]:shadow-[0_4px_16px_0_rgba(0,0,0,0.12)] data-[selected]:font-medium data-[selected]:text-gray-800 text-[13px] h-[28px]"
                          }}
                    >
                        <Tab key="download" title={
                            <div className="flex flex-row justify-center items-center ">
                                <LuDownload size={15} className='mr-2'/>
                                <span className="">Download</span>
                            </div>
                        } className="w-full">
                            <div className="-my-1 -mb-2">
                                <DownloadComponent/>
                            </div>
                        </Tab>
                        <Tab key="upload" title={
                            <div className="flex flex-row justify-center items-center ">
                                <LuUpload size={15} className='mr-2'/>
                                <span className="">Upload</span>
                            </div>
                        } className="w-full">
                            <div className="-my-1 -mb-2">
                                <UploadComponent/>
                            </div>
                        </Tab>
                    </Tabs>
                </CardBody>
            </Card>
            </FileProvider>
        </div>
    );
};

export default SharingPage;
