import React, {useState} from 'react'

import { Input, DatePicker, Snippet, Select, SelectItem, Button, Tooltip, DateValue } from '@heroui/react'
import { FaLink } from "react-icons/fa6";
import { FiExternalLink } from "react-icons/fi";
import { LuLink } from "react-icons/lu";
import { getLocalTimeZone, today } from "@internationalized/date";
import { IoMdInformationCircle } from "react-icons/io";


const ShareComponent = () => {
    const [folderName, setFolderName] = useState<string>("");
    const [sharingOption, setSharingOption] = useState<"r" | "w" | "rw">("rw")
    const [linkUrl, setLinkUrl] = useState<string | null>(null)
    const [date, setDate] = useState<DateValue | null>(null)
    const handleSumbit = () => {
        return undefined
    }

    const sharingTooltipGuide = "Once you create the folder, copy the link on the bottom to share w/ your friends."

    const sharingOptions = [
        {key: "r", label: "Read"},
        {key: "w", label: "Write"},
        {key: "rw", label: "Read/Write"}, 
    ]

    return (
        <div className="flex flex-col w-full h-full justify-center items-center p-">
            {/* Top Title */}
            <div className="flex flex-col justify-start items-center mt-1">
                <div className="flex justify-center items-center">
                    <span className="font-brsonoma font-normal text-gray-950 text-md mr-1">
                        Create a Sharing Folder
                    </span>
                    <Tooltip  placement="bottom" content={sharingTooltipGuide}>
                        <IoMdInformationCircle className='mb-1'/>
                    </Tooltip>
                </div>
                <span className="font-brsonoma font-light text-gray-500 text-xs">
                    Invite your friends to use your sharing folder
                </span>
            </div>

            {/* Read/Write Selection */}
            <div className="flex flex-row justify-start items-center my-4 p-3 border rounded-lg bg-wsecondary shadow-[0_4px_16px_0_rgba(0,0,0,0.04)]">
                <div className="flex w-8 h-8 rounded-full bg-indigo-100 justify-center items-center">
                    <LuLink className='text-blue-500' />
                </div>
                <div className="flex flex-col items-start text-nowrap ml-2">
                    <span className="font-brsonoma font-normal text-gray-950 text-xs">
                        Anyone with the link can
                    </span>
                    <span className="font-brsonoma font-light text-gray-500 text-[11px]">
                        The folder link is publicly viewable
                    </span>
                </div>
                <div className="flex w-[6.7rem] h-8 justify-center items-center ml-2 ">
                    <Select 
                        className="w-[6.7rem]"
                        multiple={false}
                        onSelectionChange={(keys) => {
                            const selectedKey = keys instanceof Set ? Array.from(keys)[0] : keys;
                            setSharingOption(selectedKey as "r" | "w" | "rw");
                          }}
                        selectedKeys={[sharingOption]}
                        classNames={{
                            trigger: 'flex items-center rounded-md border border-gray-300 bg-wsecondary h-9 min-h-9 shadow-sm',
                            mainWrapper: 'h-9',
                            value: 'text-[10px] font-brsonoma text-gray-900',
                            base: 'justify-center items-center',
                            listbox: "text-[8px]"
                        }}
                        listboxProps={{
                            itemClasses: {
                              title: "text-[12px]"
                            },
                          }}
                    >
                        {sharingOptions.map((option) => (
                        <SelectItem key={option.key}>{option.label}</SelectItem>
                        ))}
                    </Select>
                </div>
            </div>

            {/* Folder Name/Expiration Selection */}
            <div className="flex flex-col gap-3">
                <Input
                    isRequired
                    value={folderName}
                    onValueChange={setFolderName}
                    labelPlacement='outside'
                    label="How should the folder be called?"
                    placeholder='Enter a folder name'
                    className="w-[18rem]"
                    classNames={{
                        inputWrapper: 'rounded-md border border-gray-300 bg-wsecondary h-[2.5rem]',
                        innerWrapper: 'bg-transparent text-[11px] font-roboto font-light z-10',
                        label: "font-brsonoma text-[12px]",
                        input: 'font-brsonoma text-[12px]'
                    }}
                />

                <DatePicker 
                    isRequired
                    className="w-[18rem]"
                    selectorButtonPlacement="start"
                    label="When do you want the folder to expire?"
                    labelPlacement='outside'
                    value={date}
                    onChange={setDate}
                    minValue={today(getLocalTimeZone())}
                    classNames={{
                        inputWrapper: 'rounded-md border border-gray-300 bg-wsecondary h-[2.5rem]',
                        innerWrapper: 'bg-transparent font-roboto font-light',
                        calendarContent: "bg-wprimary rounded-xl",
                        label: "font-brsonoma text-[12px]",
                        input: 'font-brsonoma text-[12px]'
                    }}
                />
            </div>

            {/* Link Creation & Copy Snippet */}
            <div className="flex flex-row justify-center items-center w-[21.8rem] mt-6 gap-5">
                <div className="flex relative w-2/5 h-[2.6rem]">
                    <Snippet 
                        symbol="🔗"
                        variant="bordered"
                        classNames={{
                            base: "rounded-md border border-gray-300 bg-wsecondary h-[2.6rem] text-[10px]",
                            content: 'overflow-hidden whitespace-nowrap text-ellipsis',
                            pre: 'overflow-hidden whitespace-nowrap text-ellipsis',
                            copyButton: "text-[11px]"
                        }}
                        className="w-full"
                    >
                        {linkUrl || "Create Link"}
                    </Snippet>


                    {!linkUrl && (
                        <div className="absolute inset-0 bg-gray-300 bg-opacity-50 rounded-md flex items-center justify-center ">
                            <div className="w-full h-full bg-gradient-to-r from-gray-200 to-gray-100 bg-size-200 bg-pos-0 animate-stripes rounded-md"></div>
                        </div>
                    )}
                </div>
                <div className="flex w-3/5">
                    <Button 
                        isDisabled={linkUrl !== null}
                        color="primary"
                        className="text-xs bg-primary-gradient rounded-md w-full"
                        size="md"
                    > 
                        Create Folder
                    </Button>
                </div>
            </div>

        </div>
    )
}

export default ShareComponent