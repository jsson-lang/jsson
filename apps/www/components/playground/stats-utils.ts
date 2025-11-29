export interface StatsData {
    input: {
        lines: number;
        chars: number;
        tokens: number;
    };
    output: {
        lines: number;
        chars: number;
        tokens: number;
    };
    expansion: {
        lines: number;
        chars: number;
        tokens: number;
    };
    ratios: {
        outputPercentage: number;
        inputPercentage: number;
    };
}

export function approximateTokens(text: string): number {
    if (!text) return 0;
    return Math.ceil(text.length / 3.5);
}

export function calculateStats(input: string, output: string): StatsData {
    const inputLines = input ? input.split("\n").length : 0;
    const inputChars = input ? input.length : 0;
    const inputTokens = approximateTokens(input);

    const outputLines = output ? output.split("\n").length : 0;
    const outputChars = output ? output.length : 0;
    const outputTokens = approximateTokens(output);

    const lineExpansion = inputLines > 0 ? outputLines / inputLines : 0;
    const charExpansion = inputChars > 0 ? outputChars / inputChars : 0;
    const tokenExpansion = inputTokens > 0 ? outputTokens / inputTokens : 0;

    const totalChars = inputChars + outputChars;
    const outputPercentage = totalChars > 0 ? (outputChars / totalChars) * 100 : 0;
    const inputPercentage = totalChars > 0 ? (inputChars / totalChars) * 100 : 0;

    return {
        input: {
            lines: inputLines,
            chars: inputChars,
            tokens: inputTokens,
        },
        output: {
            lines: outputLines,
            chars: outputChars,
            tokens: outputTokens,
        },
        expansion: {
            lines: lineExpansion,
            chars: charExpansion,
            tokens: tokenExpansion,
        },
        ratios: {
            outputPercentage,
            inputPercentage,
        },
    };
}
