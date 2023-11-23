<section>
	<p>Log in to continue</p>
	<form method="POST" action="?/login">
		<label>
			Email
			<input name="email" type="email" />
		</label>
		<label>
			Password
			<input name="password" type="password" />
		</label>
		<button>Submit</button>
	</form>
</section>

<style lang="scss">
	section {
		width: 30rem;
		margin: 0 auto;
		& p {
			margin: 0 0 1rem 0;
		}
		& form {
			width: 18rem;

			display: flex;
			flex-flow: column nowrap;
			& label {
				display: flex;
				justify-content: space-between;

				& input {
					margin: 0 0 0.5rem 0;

					background-color: black;
					color: var(--primary-light);
					border: 2px solid var(--primary-light-250);
					border-radius: 5px;

					font-size: 0.85rem;
				}
			}
			& button {
				width: 30%;
				margin: 0 auto;

				background-color: transparent;
				color: var(--primary-light);
				border: 2px solid var(--primary-light);
				border-radius: 5px;

				font-weight: 500;
			}
			& button:hover {
				background-color: var(--primary-light);
				color: var(--primary-dark);

				cursor: pointer;
			}
		}
	}
</style>
