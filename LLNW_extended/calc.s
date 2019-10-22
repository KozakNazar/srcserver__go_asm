	.text
	.globl	_calc
_calc:
LFB0:
	.cfi_startproc
	pushl	%ebp
	movl	%esp, %ebp
	subl	$28, %esp
	movl	8(%ebp), %eax
	movl	%eax, -8(%ebp)
	movl	12(%ebp), %eax
	movl	%eax, -4(%ebp)
	movl	20(%ebp), %eax
	movl	%eax, -16(%ebp)
	movl	24(%ebp), %eax
	movl	%eax, -12(%ebp)
	movl	32(%ebp), %eax
	movl	%eax, -24(%ebp)
	movl	36(%ebp), %eax
	movl	%eax, -20(%ebp)
	fldl	-8(%ebp)
	fldt	LC0
	faddp	%st, %st(1)
	fldl	-16(%ebp)
	flds	16(%ebp)
	fdivrp	%st, %st(1)
	fsubrp	%st, %st(1)
	flds	28(%ebp)
	fldl	-24(%ebp)
	fmulp	%st, %st(1)
	faddp	%st, %st(1)
	fstps	-28(%ebp)
	flds	-28(%ebp)
	leave
	ret
	.cfi_endproc
LFE0:
	.section .rdata,"dr"
	.align 16
LC0:
	.long	0
	.long	-1785987072
	.long	16400
