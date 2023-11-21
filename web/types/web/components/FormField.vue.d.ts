declare const _default: __VLS_WithTemplateSlots<import("vue").DefineComponent<{
    label: {
        type: StringConstructor;
        default: null;
    };
    modelValue: {
        type: StringConstructor;
        default: null;
    };
    type: {
        type: StringConstructor;
        default: string;
    };
}, {
    focus: () => void;
}, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
    "update:modelValue": (...args: any[]) => void;
}, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
    label: {
        type: StringConstructor;
        default: null;
    };
    modelValue: {
        type: StringConstructor;
        default: null;
    };
    type: {
        type: StringConstructor;
        default: string;
    };
}>> & {
    "onUpdate:modelValue"?: ((...args: any[]) => any) | undefined;
}, {
    type: string;
    label: string;
    modelValue: string;
}, {}>, {
    default?(_: {}): any;
}>;
export default _default;
type __VLS_WithTemplateSlots<T, S> = T & {
    new (): {
        $slots: S;
    };
};
