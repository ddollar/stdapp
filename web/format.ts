import { strftime } from "./strftime";

export function day(date: string | number | Date) {
	return strftime("%Y-%m-%d", new Date(date));
}

export function firstName(name: string) {
	return (name || "").split(" ")[0];
}

export function money(amount: number | bigint) {
	if (Object.is(amount, -0)) amount = 0;
	return Intl.NumberFormat("en-US", { signDisplay: "auto", style: "currency", currency: "USD" }).format(amount);
}

export function proper(name: string) {
	return name[0].toUpperCase() + name.slice(1);
}
