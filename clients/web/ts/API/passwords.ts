import zxcvbn from "zxcvbn";

const commonWords = ["aurum", "finitum"];

export const verifyPassword = async (password: string, userInput: string[]): Promise<zxcvbn.ZXCVBNResult | string> => {
    const res = zxcvbn(password, commonWords.concat(...userInput));

    if (res.score < 2) {
        return res;
    }

    if (password.length < 8) {
        return "Password too short";
    }

    if (password.length > 72) {
        return "Password too long";
    }

    return res;
};
