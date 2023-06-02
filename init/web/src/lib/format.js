import { strftime } from "./strftime.js";

export function day(date) {
	return strftime("%Y-%m-%d", new Date(date));
}

export function money(amount) {
	if (amount === -0) amount = 0;
	return Intl.NumberFormat("en-US", { signDisplay: "auto", style: "currency", currency: "USD" }).format(amount);
}

export function proper(name) {
	return name[0].toUpperCase() + name.slice(1);
}
