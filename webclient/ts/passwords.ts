import zxcvbn from "zxcvbn";

const commonWords = ["aurum", "finitum"];

export const verifyPassword = async (password: string, userInput: string[]): Promise<zxcvbn.ZXCVBNResult> => {
    return zxcvbn(password, commonWords.concat(...userInput));
};
